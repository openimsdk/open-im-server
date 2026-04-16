package crypto

import (
	"context"
	"fmt"
	"time"

	"github.com/VirgilSecurity/virgil-sdk-go/cryptoimpl"
	"github.com/VirgilSecurity/virgil-sdk-go/sdk"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbcrypto "github.com/openimsdk/protocol/crypto"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"google.golang.org/grpc"
)

const virgilJWTTTL = 20 * time.Minute

type Config struct {
	RpcConfig     config.Crypto
	MongodbConfig config.Mongo
	Share         config.Share
	Discovery     config.Discovery
}

type cryptoServer struct {
	pbcrypto.UnimplementedCryptoServiceServer
	config       *Config
	cryptoDB     controller.CryptoDatabase
	jwtGenerator *sdk.JwtGenerator
}

func Start(ctx context.Context, config *Config, _ discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	db := mgocli.GetDB()

	deviceDB, err := mgo.NewCryptoDeviceMongo(db)
	if err != nil {
		return err
	}
	keyVersionDB, err := mgo.NewGroupKeyVersionMongo(db)
	if err != nil {
		return err
	}
	keyEventDB, err := mgo.NewGroupKeyEventMongo(db)
	if err != nil {
		return err
	}

	cryptoDB := controller.NewCryptoDatabase(deviceDB, keyVersionDB, keyEventDB, mgocli.GetTx())

	var jwtGen *sdk.JwtGenerator
	vc := config.RpcConfig.Virgil
	if vc.AppID != "" && vc.AppKey != "" && vc.AppKeyID != "" {
		virgilCrypto := cryptoimpl.NewVirgilCrypto()
		privateKey, err := virgilCrypto.ImportPrivateKey([]byte(vc.AppKey), "")
		if err != nil {
			return fmt.Errorf("import virgil app key: %w", err)
		}
		jwtGen = sdk.NewJwtGenerator(
			privateKey,
			vc.AppKeyID,
			cryptoimpl.NewVirgilAccessTokenSigner(),
			vc.AppID,
			virgilJWTTTL,
		)
	}

	pbcrypto.RegisterCryptoServiceServer(server, &cryptoServer{
		config:       config,
		cryptoDB:     cryptoDB,
		jwtGenerator: jwtGen,
	})
	return nil
}

func (s *cryptoServer) RegisterDevice(ctx context.Context, req *pbcrypto.RegisterDeviceReq) (*pbcrypto.RegisterDeviceResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "RegisterDevice request",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
		"platform", req.Platform,
		"deviceModel", req.DeviceModel,
		"appVersion", req.AppVersion,
	)
	if req.UserID == "" || req.DeviceID == "" {
		log.ZError(ctx, "RegisterDevice invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, errs.ErrArgs.WrapMsg("userID and deviceID are required")
	}
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		log.ZError(ctx, "RegisterDevice auth check failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, err
	}
	device, err := s.cryptoDB.RegisterDevice(ctx, req.UserID, req.DeviceID, req.Platform, req.DeviceModel, req.AppVersion)
	if err != nil {
		log.ZError(ctx, "RegisterDevice db failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"platform", req.Platform,
		)
		return nil, err
	}
	log.ZDebug(ctx, "RegisterDevice success",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
		"virgilIdentity", device.VirgilIdentity,
	)
	return &pbcrypto.RegisterDeviceResp{
		Device: modelDeviceToProto(device),
	}, nil
}

func (s *cryptoServer) GetDevices(ctx context.Context, req *pbcrypto.GetDevicesReq) (*pbcrypto.GetDevicesResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "GetDevices request", "opUserID", opUserID, "targetUserID", req.UserID)
	if req.UserID == "" {
		log.ZError(ctx, "GetDevices invalid args", errs.ErrArgs, "opUserID", opUserID, "targetUserID", req.UserID)
		return nil, errs.ErrArgs.WrapMsg("userID is required")
	}
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		log.ZError(ctx, "GetDevices auth check failed", err, "opUserID", opUserID, "targetUserID", req.UserID)
		return nil, err
	}
	devices, err := s.cryptoDB.GetDevices(ctx, req.UserID)
	if err != nil {
		log.ZError(ctx, "GetDevices db failed", err, "opUserID", opUserID, "targetUserID", req.UserID)
		return nil, err
	}
	pbDevices := make([]*pbcrypto.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		pbDevices = append(pbDevices, modelDeviceToProto(d))
	}
	log.ZDebug(ctx, "GetDevices success",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceCount", len(devices),
	)
	return &pbcrypto.GetDevicesResp{Devices: pbDevices}, nil
}

func (s *cryptoServer) RevokeDevice(ctx context.Context, req *pbcrypto.RevokeDeviceReq) (*pbcrypto.RevokeDeviceResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "RevokeDevice request",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
	)
	if req.UserID == "" || req.DeviceID == "" {
		log.ZError(ctx, "RevokeDevice invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, errs.ErrArgs.WrapMsg("userID and deviceID are required")
	}
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		log.ZError(ctx, "RevokeDevice auth check failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, err
	}
	if err := s.cryptoDB.RevokeDevice(ctx, req.UserID, req.DeviceID); err != nil {
		log.ZError(ctx, "RevokeDevice db failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, err
	}
	log.ZDebug(ctx, "RevokeDevice success",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
	)
	return &pbcrypto.RevokeDeviceResp{}, nil
}

func (s *cryptoServer) GetVirgilJWT(ctx context.Context, req *pbcrypto.GetVirgilJWTReq) (*pbcrypto.GetVirgilJWTResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "GetVirgilJWT request",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
	)
	if req.UserID == "" || req.DeviceID == "" {
		log.ZError(ctx, "GetVirgilJWT invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, errs.ErrArgs.WrapMsg("userID and deviceID are required")
	}
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		log.ZError(ctx, "GetVirgilJWT auth check failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, err
	}
	if s.jwtGenerator == nil {
		log.ZError(ctx, "GetVirgilJWT jwt generator not configured", errs.New("virgil is not configured"),
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, errs.New("virgil is not configured").Wrap()
	}

	device, err := s.cryptoDB.GetDevice(ctx, req.UserID, req.DeviceID)
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			log.ZError(ctx, "GetVirgilJWT device not found", err,
				"opUserID", opUserID,
				"targetUserID", req.UserID,
				"deviceID", req.DeviceID,
			)
			return nil, errs.ErrRecordNotFound.WrapMsg("device not found")
		}
		log.ZError(ctx, "GetVirgilJWT query device failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
		)
		return nil, err
	}
	if device.Status != "active" {
		log.ZError(ctx, "GetVirgilJWT device revoked", errs.ErrNoPermission,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"status", device.Status,
		)
		return nil, errs.ErrNoPermission.WrapMsg("device is revoked")
	}

	if err := s.cryptoDB.TouchDevice(ctx, req.UserID, req.DeviceID); err != nil {
		log.ZError(ctx, "TouchDevice failed", err,
			"opUserID", opUserID,
			"userID", req.UserID,
			"deviceID", req.DeviceID,
		)
	}

	identity := req.UserID + ":" + req.DeviceID
	token, err := s.jwtGenerator.GenerateToken(identity, nil)
	if err != nil {
		log.ZError(ctx, "GetVirgilJWT generate token failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"virgilIdentity", identity,
		)
		return nil, errs.New("generate virgil jwt failed").Wrap()
	}
	log.ZDebug(ctx, "GetVirgilJWT success",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
		"virgilIdentity", identity,
		"expiresInSec", int64(virgilJWTTTL/time.Second),
	)

	return &pbcrypto.GetVirgilJWTResp{
		VirgilJWT:      token.String(),
		ExpiresIn:      int64(virgilJWTTTL / time.Second),
		VirgilIdentity: identity,
	}, nil
}

func (s *cryptoServer) GetGroupKeyVersion(ctx context.Context, req *pbcrypto.GetGroupKeyVersionReq) (*pbcrypto.GetGroupKeyVersionResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "GetGroupKeyVersion request", "opUserID", opUserID, "groupID", req.GroupID)
	if req.GroupID == "" {
		log.ZError(ctx, "GetGroupKeyVersion invalid args", errs.ErrArgs, "opUserID", opUserID, "groupID", req.GroupID)
		return nil, errs.ErrArgs.WrapMsg("groupID is required")
	}
	version, err := s.cryptoDB.GetGroupKeyVersion(ctx, req.GroupID)
	if err != nil {
		log.ZError(ctx, "GetGroupKeyVersion db failed", err, "opUserID", opUserID, "groupID", req.GroupID)
		return nil, err
	}
	log.ZDebug(ctx, "GetGroupKeyVersion success",
		"opUserID", opUserID,
		"groupID", req.GroupID,
		"groupKeyVersion", version,
	)
	return &pbcrypto.GetGroupKeyVersionResp{
		GroupID:         req.GroupID,
		GroupKeyVersion: version,
	}, nil
}

// BumpGroupKeyVersion is internal-only (not exposed via HTTP API).
// Called by Group Service after membership changes.
func (s *cryptoServer) BumpGroupKeyVersion(ctx context.Context, req *pbcrypto.BumpGroupKeyVersionReq) (*pbcrypto.BumpGroupKeyVersionResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "BumpGroupKeyVersion request",
		"opUserID", opUserID,
		"groupID", req.GroupID,
		"operatorUserID", req.OperatorUserID,
		"eventType", req.EventType,
	)
	if req.GroupID == "" {
		log.ZError(ctx, "BumpGroupKeyVersion invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"groupID", req.GroupID,
			"operatorUserID", req.OperatorUserID,
			"eventType", req.EventType,
		)
		return nil, errs.ErrArgs.WrapMsg("groupID is required")
	}
	newVersion, err := s.cryptoDB.BumpGroupKeyVersion(ctx, req.GroupID, req.OperatorUserID, req.EventType)
	if err != nil {
		log.ZError(ctx, "BumpGroupKeyVersion db failed", err,
			"opUserID", opUserID,
			"groupID", req.GroupID,
			"operatorUserID", req.OperatorUserID,
			"eventType", req.EventType,
		)
		return nil, err
	}
	log.ZDebug(ctx, "BumpGroupKeyVersion success",
		"opUserID", opUserID,
		"groupID", req.GroupID,
		"operatorUserID", req.OperatorUserID,
		"eventType", req.EventType,
		"newGroupKeyVersion", newVersion,
	)
	return &pbcrypto.BumpGroupKeyVersionResp{
		GroupID:         req.GroupID,
		GroupKeyVersion: newVersion,
	}, nil
}

func (s *cryptoServer) GetGroupKeyEvents(ctx context.Context, req *pbcrypto.GetGroupKeyEventsReq) (*pbcrypto.GetGroupKeyEventsResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "GetGroupKeyEvents request",
		"opUserID", opUserID,
		"groupID", req.GroupID,
		"sinceVersion", req.SinceVersion,
	)
	if req.GroupID == "" {
		log.ZError(ctx, "GetGroupKeyEvents invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"groupID", req.GroupID,
			"sinceVersion", req.SinceVersion,
		)
		return nil, errs.ErrArgs.WrapMsg("groupID is required")
	}
	events, err := s.cryptoDB.GetGroupKeyEvents(ctx, req.GroupID, req.SinceVersion)
	if err != nil {
		log.ZError(ctx, "GetGroupKeyEvents db failed", err,
			"opUserID", opUserID,
			"groupID", req.GroupID,
			"sinceVersion", req.SinceVersion,
		)
		return nil, err
	}
	pbEvents := make([]*pbcrypto.GroupKeyEventInfo, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, &pbcrypto.GroupKeyEventInfo{
			EventID:         e.EventID,
			GroupID:         e.GroupID,
			GroupKeyVersion: e.GroupKeyVersion,
			EventType:       e.EventType,
			OperatorUserID:  e.OperatorUserID,
			CreateTime:      e.CreateTime.UnixMilli(),
		})
	}
	log.ZDebug(ctx, "GetGroupKeyEvents success",
		"opUserID", opUserID,
		"groupID", req.GroupID,
		"sinceVersion", req.SinceVersion,
		"eventCount", len(events),
	)
	return &pbcrypto.GetGroupKeyEventsResp{Events: pbEvents}, nil
}

func (s *cryptoServer) SecurityPrecheck(ctx context.Context, req *pbcrypto.SecurityPrecheckReq) (*pbcrypto.SecurityPrecheckResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "SecurityPrecheck request",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
		"action", req.Action,
	)
	if req.UserID == "" || req.DeviceID == "" {
		log.ZError(ctx, "SecurityPrecheck invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"action", req.Action,
		)
		return nil, errs.ErrArgs.WrapMsg("userID and deviceID are required")
	}
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		log.ZError(ctx, "SecurityPrecheck auth check failed", err,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"action", req.Action,
		)
		return nil, err
	}
	device, err := s.cryptoDB.GetDevice(ctx, req.UserID, req.DeviceID)
	if err != nil {
		log.ZDebug(ctx, "SecurityPrecheck denied",
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"reason", "device not found",
		)
		return &pbcrypto.SecurityPrecheckResp{Allowed: false, Reason: "device not found"}, nil
	}
	if device.Status != "active" {
		log.ZDebug(ctx, "SecurityPrecheck denied",
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"reason", "device is revoked",
		)
		return &pbcrypto.SecurityPrecheckResp{Allowed: false, Reason: "device is revoked"}, nil
	}
	log.ZDebug(ctx, "SecurityPrecheck allowed",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
	)
	return &pbcrypto.SecurityPrecheckResp{Allowed: true}, nil
}

// IntegrityReport is a placeholder for future device integrity verification.
// Currently accepts all reports; implement validation logic when ready.
func (s *cryptoServer) IntegrityReport(ctx context.Context, req *pbcrypto.IntegrityReportReq) (*pbcrypto.IntegrityReportResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	log.ZDebug(ctx, "IntegrityReport request",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
		"timestamp", req.Timestamp,
		"reportSize", len(req.ReportData),
	)
	if req.UserID == "" || req.DeviceID == "" {
		log.ZError(ctx, "IntegrityReport invalid args", errs.ErrArgs,
			"opUserID", opUserID,
			"targetUserID", req.UserID,
			"deviceID", req.DeviceID,
			"timestamp", req.Timestamp,
		)
		return nil, errs.ErrArgs.WrapMsg("userID and deviceID are required")
	}
	log.ZDebug(ctx, "IntegrityReport accepted",
		"opUserID", opUserID,
		"targetUserID", req.UserID,
		"deviceID", req.DeviceID,
	)
	return &pbcrypto.IntegrityReportResp{Accepted: true}, nil
}

func modelDeviceToProto(d *model.CryptoDevice) *pbcrypto.DeviceInfo {
	return &pbcrypto.DeviceInfo{
		DeviceID:       d.DeviceID,
		UserID:         d.UserID,
		Platform:       d.Platform,
		DeviceModel:    d.DeviceModel,
		AppVersion:     d.AppVersion,
		VirgilIdentity: d.VirgilIdentity,
		Status:         d.Status,
		LastSeenAt:     d.LastSeenAt.UnixMilli(),
		CreateTime:     d.CreateTime.UnixMilli(),
	}
}
