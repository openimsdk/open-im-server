package captcha

import (
	"context"
	"encoding/json"
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

	"github.com/wenlng/go-captcha/v2/click"
)

// alphanumChars is the character pool for the click captcha.
// Visually ambiguous characters (I, O, 0, 1, l) are excluded.
var alphanumChars = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "J", "K",
	"L", "M", "N", "P", "Q", "R", "S", "T", "U", "V",
	"W", "X", "Y", "Z", "2", "3", "4", "5", "6", "7", "8", "9",
}

type Config struct {
	RpcConfig     config.Captcha
	MongodbConfig config.Mongo
	Share         config.Share
	Discovery     config.Discovery
}

type server struct {
	pbcaptcha.UnimplementedCaptchaServer
	conf       config.Captcha
	capt       click.Captcha
	collection *mongo.Collection
}

// captchaDoc is the MongoDB document that stores the verification answer.
type captchaDoc struct {
	CaptchaID  string    `bson:"captcha_id"`
	DotsJSON   string    `bson:"dots_json"` // JSON-encoded map[int]*click.Dot (answer dots)
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

	capt, err := buildClickCaptcha()
	if err != nil {
		log.ZError(ctx, "captcha build click captcha failed", err)
		return err
	}

	s := &server{
		conf:       cfg.RpcConfig,
		capt:       capt,
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

	dots := captData.GetData() // answer dots: map[int]*click.Dot
	masterImage, err := captData.GetMasterImage().ToBase64DataWithQuality(0)
	if err != nil {
		log.ZError(ctx, "captcha encode master image failed", err)
		return nil, err
	}
	thumbImage, err := captData.GetThumbImage().ToBase64Data()
	if err != nil {
		log.ZError(ctx, "captcha encode thumb image failed", err)
		return nil, err
	}

	dotsJSON, err := json.Marshal(dots)
	if err != nil {
		log.ZError(ctx, "captcha marshal dots failed", err)
		return nil, err
	}

	id := uuid.NewString()
	now := time.Now()
	expiredAt := now.Add(time.Duration(s.conf.ExpireSeconds) * time.Second)
	_, err = s.collection.InsertOne(ctx, captchaDoc{
		CaptchaID:  id,
		DotsJSON:   string(dotsJSON),
		ExpiredAt:  expiredAt,
		CreateTime: now,
	})
	if err != nil {
		log.ZError(ctx, "captcha insert mongodb failed", err, "captchaID", id)
		return nil, err
	}

	log.ZDebug(ctx, "captcha generated", "captchaID", id, "dotCount", len(dots), "expireAt", expiredAt.Unix())
	return &pbcaptcha.GenerateCaptchaResp{
		CaptchaID:   id,
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
		ExpireAt:    expiredAt.Unix(),
	}, nil
}

func (s *server) VerifyCaptcha(ctx context.Context, req *pbcaptcha.VerifyCaptchaReq) (*pbcaptcha.VerifyCaptchaResp, error) {
	log.ZDebug(ctx, "captcha verify request", "captchaID", req.CaptchaID, "clickCount", len(req.ClickPoints))

	now := time.Now()
	filter := bson.M{
		"captcha_id":  req.CaptchaID,
		"verify_time": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"verify_time": now}}

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

	// Unmarshal the stored answer dots.
	var answerDots map[int]*click.Dot
	if err := json.Unmarshal([]byte(doc.DotsJSON), &answerDots); err != nil {
		log.ZError(ctx, "captcha unmarshal dots failed", err, "captchaID", req.CaptchaID)
		return nil, servererrs.ErrDatabase.WrapMsg("internal captcha data error")
	}

	success := validateClickPoints(req.ClickPoints, answerDots, s.conf.VerifyPadding)
	if !success {
		log.ZError(ctx, "captcha validate failed", nil,
			"captchaID", req.CaptchaID,
			"clickCount", len(req.ClickPoints),
			"answerCount", len(answerDots),
		)
	} else {
		log.ZDebug(ctx, "captcha validate success", "captchaID", req.CaptchaID)
	}
	return &pbcaptcha.VerifyCaptchaResp{Success: success}, nil
}

// validateClickPoints checks that each user click point falls within the
// bounding box of the corresponding answer dot (in order).
func validateClickPoints(points []*pbcaptcha.ClickPoint, dots map[int]*click.Dot, padding int) bool {
	if len(points) != len(dots) {
		return false
	}
	for i, pt := range points {
		dot, ok := dots[i]
		if !ok {
			return false
		}
		if !click.Validate(int(pt.X), int(pt.Y), dot.X, dot.Y, dot.Width, dot.Height, padding) {
			return false
		}
	}
	return true
}
