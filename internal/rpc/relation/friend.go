package relation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/notification/common_user"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"

	"github.com/openimsdk/tools/mq/memamq"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
)

type friendServer struct {
	relation.UnimplementedFriendServer
	db                 controller.FriendDatabase
	blackDatabase      controller.BlackDatabase
	notificationSender *FriendNotificationSender
	RegisterCenter     discovery.Conn
	config             *Config
	webhookClient      *webhook.Client
	queue              *memamq.MemoryQueue
	userClient         *rpcli.UserClient
}

type Config struct {
	RpcConfig     config.Friend
	RedisConfig   config.Redis
	MongodbConfig config.Mongo
	// ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error {
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}

	friendMongoDB, err := mgo.NewFriendMongo(mgocli.GetDB())
	if err != nil {
		return err
	}

	friendRequestMongoDB, err := mgo.NewFriendRequestMongo(mgocli.GetDB())
	if err != nil {
		return err
	}

	blackMongoDB, err := mgo.NewBlackMongo(mgocli.GetDB())
	if err != nil {
		return err
	}

	userConn, err := client.GetConn(ctx, config.Discovery.RpcService.User)
	if err != nil {
		return err
	}
	msgConn, err := client.GetConn(ctx, config.Discovery.RpcService.Msg)
	if err != nil {
		return err
	}
	userClient := rpcli.NewUserClient(userConn)
	database := controller.NewFriendDatabase(
		friendMongoDB,
		friendRequestMongoDB,
		redis.NewFriendCacheRedis(rdb, &config.LocalCacheConfig, friendMongoDB),
		mgocli.GetTx(),
	)
	// Initialize notification sender
	notificationSender := NewFriendNotificationSender(
		&config.NotificationConfig,
		rpcli.NewMsgClient(msgConn),
		WithRpcFunc(userClient.GetUsersInfo),
		WithFriendDB(database),
	)
	localcache.InitLocalCache(&config.LocalCacheConfig)

	// Register Friend server with refactored MongoDB and Redis integrations
	relation.RegisterFriendServer(server, &friendServer{
		db: database,
		blackDatabase: controller.NewBlackDatabase(
			blackMongoDB,
			redis.NewBlackCacheRedis(rdb, &config.LocalCacheConfig, blackMongoDB),
		),
		notificationSender: notificationSender,
		RegisterCenter:     client,
		config:             config,
		webhookClient:      webhook.NewWebhookClient(config.WebhooksConfig.URL),
		queue:              memamq.NewMemoryQueue(16, 1024*1024),
		userClient:         userClient,
	})
	return nil
}

// ok.
func (s *friendServer) ApplyToAddFriend(ctx context.Context, req *relation.ApplyToAddFriendReq) (resp *relation.ApplyToAddFriendResp, err error) {
	resp = &relation.ApplyToAddFriendResp{}
	if err := authverify.CheckAccess(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if req.ToUserID == req.FromUserID {
		return nil, servererrs.ErrCanNotAddYourself.WrapMsg("req.ToUserID", req.ToUserID)
	}
	if err = s.webhookBeforeAddFriend(ctx, &s.config.WebhooksConfig.BeforeAddFriend, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}
	if err := s.userClient.CheckUser(ctx, []string{req.ToUserID, req.FromUserID}); err != nil {
		return nil, err
	}

	in1, in2, err := s.db.CheckIn(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, servererrs.ErrRelationshipAlready.WrapMsg("already friends has f")
	}
	if err = s.db.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	s.notificationSender.FriendApplicationAddNotification(ctx, req)
	s.webhookAfterAddFriend(ctx, &s.config.WebhooksConfig.AfterAddFriend, req)
	return resp, nil
}

// ok.
func (s *friendServer) ImportFriends(ctx context.Context, req *relation.ImportFriendReq) (resp *relation.ImportFriendResp, err error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}

	if err := s.userClient.CheckUser(ctx, append([]string{req.OwnerUserID}, req.FriendUserIDs...)); err != nil {
		return nil, err
	}
	if datautil.Contain(req.OwnerUserID, req.FriendUserIDs...) {
		return nil, servererrs.ErrCanNotAddYourself.WrapMsg("can not add yourself")
	}
	if datautil.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friend userID repeated")
	}

	if err := s.webhookBeforeImportFriends(ctx, &s.config.WebhooksConfig.BeforeImportFriends, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	if err := s.db.BecomeFriends(ctx, req.OwnerUserID, req.FriendUserIDs, constant.BecomeFriendByImport); err != nil {
		return nil, err
	}
	for _, userID := range req.FriendUserIDs {
		s.notificationSender.FriendApplicationAgreedNotification(ctx, &relation.RespondFriendApplyReq{
			FromUserID:   req.OwnerUserID,
			ToUserID:     userID,
			HandleResult: constant.FriendResponseAgree,
		}, false)
	}

	s.webhookAfterImportFriends(ctx, &s.config.WebhooksConfig.AfterImportFriends, req)
	return &relation.ImportFriendResp{}, nil
}

// ok.
func (s *friendServer) RespondFriendApply(ctx context.Context, req *relation.RespondFriendApplyReq) (resp *relation.RespondFriendApplyResp, err error) {
	resp = &relation.RespondFriendApplyResp{}
	if err := authverify.CheckAccess(ctx, req.ToUserID); err != nil {
		return nil, err
	}

	friendRequest := model.FriendRequest{
		FromUserID:   req.FromUserID,
		ToUserID:     req.ToUserID,
		HandleMsg:    req.HandleMsg,
		HandleResult: req.HandleResult,
	}
	if req.HandleResult == constant.FriendResponseAgree {
		if err := s.webhookBeforeAddFriendAgree(ctx, &s.config.WebhooksConfig.BeforeAddFriendAgree, req); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}
		err := s.db.AgreeFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.webhookAfterAddFriendAgree(ctx, &s.config.WebhooksConfig.AfterAddFriendAgree, req)
		s.notificationSender.FriendApplicationAgreedNotification(ctx, req, true)
		return resp, nil
	}
	if req.HandleResult == constant.FriendResponseRefuse {
		err := s.db.RefuseFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.notificationSender.FriendApplicationRefusedNotification(ctx, req)
		return resp, nil
	}
	return nil, errs.ErrArgs.WrapMsg("req.HandleResult != -1/1")
}

// ok.
func (s *friendServer) DeleteFriend(ctx context.Context, req *relation.DeleteFriendReq) (resp *relation.DeleteFriendResp, err error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	_, err = s.db.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}

	if err := s.db.Delete(ctx, req.OwnerUserID, []string{req.FriendUserID}); err != nil {
		return nil, err
	}

	s.notificationSender.FriendDeletedNotification(ctx, req)
	s.webhookAfterDeleteFriend(ctx, &s.config.WebhooksConfig.AfterDeleteFriend, req)

	return &relation.DeleteFriendResp{}, nil
}

// ok.
func (s *friendServer) SetFriendRemark(ctx context.Context, req *relation.SetFriendRemarkReq) (resp *relation.SetFriendRemarkResp, err error) {
	if err = s.webhookBeforeSetFriendRemark(ctx, &s.config.WebhooksConfig.BeforeSetFriendRemark, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	_, err = s.db.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}

	if err := s.db.UpdateRemark(ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}

	s.webhookAfterSetFriendRemark(ctx, &s.config.WebhooksConfig.AfterSetFriendRemark, req)
	s.notificationSender.FriendRemarkSetNotification(ctx, req.OwnerUserID, req.FriendUserID)

	return &relation.SetFriendRemarkResp{}, nil
}

func (s *friendServer) GetFriendInfo(ctx context.Context, req *relation.GetFriendInfoReq) (*relation.GetFriendInfoResp, error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	friends, err := s.db.FindFriendsWithError(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	return &relation.GetFriendInfoResp{FriendInfos: convert.FriendOnlyDB2PbOnly(friends)}, nil
}

func (s *friendServer) GetDesignatedFriends(ctx context.Context, req *relation.GetDesignatedFriendsReq) (resp *relation.GetDesignatedFriendsResp, err error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	resp = &relation.GetDesignatedFriendsResp{}
	if datautil.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friend userID repeated")
	}
	friends, err := s.getFriend(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	return &relation.GetDesignatedFriendsResp{
		FriendsInfo: friends,
	}, nil
}

func (s *friendServer) getFriend(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*sdkws.FriendInfo, error) {
	if len(friendUserIDs) == 0 {
		return nil, nil
	}
	friends, err := s.db.FindFriendsWithError(ctx, ownerUserID, friendUserIDs)
	if err != nil {
		return nil, err
	}
	return convert.FriendsDB2Pb(ctx, friends, s.userClient.GetUsersInfoMap)
}

// Get the list of friend requests sent out proactively.
func (s *friendServer) GetDesignatedFriendsApply(ctx context.Context, req *relation.GetDesignatedFriendsApplyReq) (resp *relation.GetDesignatedFriendsApplyResp, err error) {
	if err := authverify.CheckAccessIn(ctx, req.FromUserID, req.ToUserID); err != nil {
		return nil, err
	}
	friendRequests, err := s.db.FindBothFriendRequests(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	resp = &relation.GetDesignatedFriendsApplyResp{}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.getCommonUserMap)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get received friend requests (i.e., those initiated by others).
func (s *friendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *relation.GetPaginationFriendsApplyToReq) (resp *relation.GetPaginationFriendsApplyToResp, err error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}

	handleResults := datautil.Slice(req.HandleResults, func(e int32) int {
		return int(e)
	})
	total, friendRequests, err := s.db.PageFriendRequestToMe(ctx, req.UserID, handleResults, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp = &relation.GetPaginationFriendsApplyToResp{}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.getCommonUserMap)
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)

	return resp, nil
}

func (s *friendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *relation.GetPaginationFriendsApplyFromReq) (resp *relation.GetPaginationFriendsApplyFromResp, err error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}

	handleResults := datautil.Slice(req.HandleResults, func(e int32) int {
		return int(e)
	})
	total, friendRequests, err := s.db.PageFriendRequestFromMe(ctx, req.UserID, handleResults, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp = &relation.GetPaginationFriendsApplyFromResp{}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.getCommonUserMap)
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)

	return resp, nil
}

// ok.
func (s *friendServer) IsFriend(ctx context.Context, req *relation.IsFriendReq) (resp *relation.IsFriendResp, err error) {
	if err := authverify.CheckAccessIn(ctx, req.UserID1, req.UserID2); err != nil {
		return nil, err
	}
	resp = &relation.IsFriendResp{}
	resp.InUser1Friends, resp.InUser2Friends, err = s.db.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *friendServer) GetPaginationFriends(ctx context.Context, req *relation.GetPaginationFriendsReq) (resp *relation.GetPaginationFriendsResp, err error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}

	total, friends, err := s.db.PageOwnerFriends(ctx, req.UserID, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp = &relation.GetPaginationFriendsResp{}
	resp.FriendsInfo, err = convert.FriendsDB2Pb(ctx, friends, s.userClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)

	return resp, nil
}

func (s *friendServer) GetFriendIDs(ctx context.Context, req *relation.GetFriendIDsReq) (resp *relation.GetFriendIDsResp, err error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}

	resp = &relation.GetFriendIDsResp{}
	resp.FriendIDs, err = s.db.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *friendServer) GetSpecifiedFriendsInfo(ctx context.Context, req *relation.GetSpecifiedFriendsInfoReq) (*relation.GetSpecifiedFriendsInfoResp, error) {
	if len(req.UserIDList) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDList is empty")
	}

	if datautil.Duplicate(req.UserIDList) {
		return nil, errs.ErrArgs.WrapMsg("userIDList repeated")
	}

	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	userMap, err := s.userClient.GetUsersInfoMap(ctx, req.UserIDList)
	if err != nil {
		return nil, err
	}

	friends, err := s.db.FindFriendsWithError(ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}

	blacks, err := s.blackDatabase.FindBlackInfos(ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}

	friendMap := datautil.SliceToMap(friends, func(e *model.Friend) string {
		return e.FriendUserID
	})

	blackMap := datautil.SliceToMap(blacks, func(e *model.Black) string {
		return e.BlockUserID
	})

	resp := &relation.GetSpecifiedFriendsInfoResp{
		Infos: make([]*relation.GetSpecifiedFriendsInfoInfo, 0, len(req.UserIDList)),
	}

	for _, userID := range req.UserIDList {
		user := userMap[userID]

		if user == nil {
			continue
		}

		var friendInfo *sdkws.FriendInfo
		if friend := friendMap[userID]; friend != nil {
			friendInfo = &sdkws.FriendInfo{
				OwnerUserID:    friend.OwnerUserID,
				Remark:         friend.Remark,
				CreateTime:     friend.CreateTime.UnixMilli(),
				AddSource:      friend.AddSource,
				OperatorUserID: friend.OperatorUserID,
				Ex:             friend.Ex,
				IsPinned:       friend.IsPinned,
			}
		}

		var blackInfo *sdkws.BlackInfo
		if black := blackMap[userID]; black != nil {
			blackInfo = &sdkws.BlackInfo{
				OwnerUserID:    black.OwnerUserID,
				CreateTime:     black.CreateTime.UnixMilli(),
				AddSource:      black.AddSource,
				OperatorUserID: black.OperatorUserID,
				Ex:             black.Ex,
			}
		}

		resp.Infos = append(resp.Infos, &relation.GetSpecifiedFriendsInfoInfo{
			UserInfo:   user,
			FriendInfo: friendInfo,
			BlackInfo:  blackInfo,
		})
	}

	return resp, nil
}

func (s *friendServer) UpdateFriends(ctx context.Context, req *relation.UpdateFriendsReq) (*relation.UpdateFriendsResp, error) {
	if len(req.FriendUserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("friendIDList is empty")
	}
	if datautil.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friendIDList repeated")
	}

	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	_, err := s.db.FindFriendsWithError(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	val := make(map[string]any)

	if req.IsPinned != nil {
		val["is_pinned"] = req.IsPinned.Value
	}
	if req.Remark != nil {
		val["remark"] = req.Remark.Value
	}
	if req.Ex != nil {
		val["ex"] = req.Ex.Value
	}
	if err = s.db.UpdateFriends(ctx, req.OwnerUserID, req.FriendUserIDs, val); err != nil {
		return nil, err
	}

	resp := &relation.UpdateFriendsResp{}

	s.notificationSender.FriendsInfoUpdateNotification(ctx, req.OwnerUserID, req.FriendUserIDs)
	return resp, nil
}

func (s *friendServer) GetSelfUnhandledApplyCount(ctx context.Context, req *relation.GetSelfUnhandledApplyCountReq) (*relation.GetSelfUnhandledApplyCountResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}

	count, err := s.db.GetUnhandledCount(ctx, req.UserID, req.Time)
	if err != nil {
		return nil, err
	}

	return &relation.GetSelfUnhandledApplyCountResp{
		Count: count,
	}, nil
}

func (s *friendServer) getCommonUserMap(ctx context.Context, userIDs []string) (map[string]common_user.CommonUser, error) {
	users, err := s.userClient.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMapAny(users, func(e *sdkws.UserInfo) (string, common_user.CommonUser) {
		return e.UserID, e
	}), nil
}
