package captcha

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	pbcaptcha "github.com/openimsdk/protocol/captcha"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Config struct {
	RpcConfig     config.Captcha
	MongodbConfig config.Mongo
	Share         config.Share
	Discovery     config.Discovery
}

type server struct {
	pbcaptcha.UnimplementedCaptchaServer
	conf       config.Captcha
	capt       slide.Captcha
	collection *mongo.Collection
}

type captchaDoc struct {
	CaptchaID  string    `bson:"captcha_id"`
	X          int       `bson:"x"`
	Y          int       `bson:"y"`
	ExpiredAt  time.Time `bson:"expired_at"`
	CreateTime time.Time `bson:"create_time"`
	VerifyTime time.Time `bson:"verify_time,omitempty"`
}

func Start(ctx context.Context, cfg *Config, _ discovery.SvcDiscoveryRegistry, grpcServer *grpc.Server) error {
	mongoClient, err := mongoutil.NewMongoDB(ctx, cfg.MongodbConfig.Build())
	if err != nil {
		log.ZError(ctx, "captcha connect mongodb failed", err)
		return err
	}
	collection := mongoClient.GetDB().Collection("captcha")
	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "captcha_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expired_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	if err != nil {
		log.ZError(ctx, "captcha create mongodb indexes failed", err)
		return err
	}

	resources, err := loadResources()
	if err != nil {
		log.ZError(ctx, "captcha load resources failed", err)
		return err
	}

	builder := slide.NewBuilder()
	builder.SetResources(resources...)
	s := &server{
		conf:       cfg.RpcConfig,
		capt:       builder.Make(),
		collection: collection,
	}
	if s.conf.ExpireSeconds <= 0 {
		s.conf.ExpireSeconds = 120
	}
	if s.conf.VerifyPadding <= 0 {
		s.conf.VerifyPadding = 8
	}
	pbcaptcha.RegisterCaptchaServer(grpcServer, s)
	return nil
}

func (s *server) GenerateCaptcha(ctx context.Context, _ *pbcaptcha.GenerateCaptchaReq) (*pbcaptcha.GenerateCaptchaResp, error) {
	captData, err := s.capt.Generate()
	if err != nil {
		log.ZError(ctx, "captcha generate failed", err)
		return nil, err
	}
	block := captData.GetData()
	masterImage, err := captData.GetMasterImage().ToBase64DataWithQuality(option.QualityNone)
	if err != nil {
		log.ZError(ctx, "captcha encode master image failed", err)
		return nil, err
	}
	tileImage, err := captData.GetTileImage().ToBase64Data()
	if err != nil {
		log.ZError(ctx, "captcha encode tile image failed", err)
		return nil, err
	}
	id := uuid.NewString()
	now := time.Now()
	expiredAt := now.Add(time.Duration(s.conf.ExpireSeconds) * time.Second)
	_, err = s.collection.InsertOne(ctx, captchaDoc{
		CaptchaID:  id,
		X:          block.X,
		Y:          block.Y,
		ExpiredAt:  expiredAt,
		CreateTime: now,
	})
	if err != nil {
		log.ZError(ctx, "captcha insert mongodb failed", err, "captchaID", id)
		return nil, err
	}
	return &pbcaptcha.GenerateCaptchaResp{
		CaptchaID:   id,
		MasterImage: masterImage,
		TileImage:   tileImage,
		ExpireAt:    expiredAt.Unix(),
		TileY:       int32(block.Y),
	}, nil
}

func (s *server) VerifyCaptcha(ctx context.Context, req *pbcaptcha.VerifyCaptchaReq) (*pbcaptcha.VerifyCaptchaResp, error) {
	now := time.Now()
	filter := bson.M{
		"captcha_id":  req.CaptchaID,
		"verify_time": bson.M{"$exists": false},
	}
	update := bson.M{
		"$set": bson.M{
			"verify_time": now,
		},
	}
	var doc captchaDoc
	err := s.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.Before),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.ZWarn(ctx, "captcha not found or already verified", err, "captchaID", req.CaptchaID)
			return nil, servererrs.ErrRecordNotFound.WrapMsg("captcha not found, expired, or already verified", "captchaID", req.CaptchaID)
		}
		log.ZError(ctx, "captcha verify query failed", err, "captchaID", req.CaptchaID)
		return nil, servererrs.ErrDatabase.WrapMsg("verify captcha query failed", "captchaID", req.CaptchaID)
	}
	if now.After(doc.ExpiredAt) {
		log.ZWarn(ctx, "captcha expired", nil, "captchaID", req.CaptchaID, "expiredAt", doc.ExpiredAt.Unix())
		return nil, servererrs.ErrFileUploadedExpired.WrapMsg("captcha expired", "captchaID", req.CaptchaID)
	}
	x, y := req.GetX(), req.GetY()
	success := slide.Validate(int(x), int(y), doc.X, doc.Y, s.conf.VerifyPadding)
	if !success {
		log.ZError(ctx, "captcha validate failed", nil, "captchaID", req.CaptchaID, "x", x, "y", y, "docX", doc.X, "docY", doc.Y)
	}
	return &pbcaptcha.VerifyCaptchaResp{Success: success}, nil
}
