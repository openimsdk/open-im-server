package office

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/office"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"google.golang.org/grpc"
)

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	//rdb, err := cache.NewRedis()
	//if err != nil {
	//	return err
	//}
	office.RegisterOfficeServer(server, &officeServer{
		officeDatabase: controller.NewOfficeDatabase(mongo),
		msgRpcClient:   rpcclient.NewMsgClient(client),
		user:           rpcclient.NewUserClient(client),
	})
	return nil
}

type officeServer struct {
	officeDatabase controller.OfficeDatabase
	user           *rpcclient.UserClient
	msgRpcClient   *rpcclient.MsgClient
}

func (o *officeServer) GetUserTags(ctx context.Context, req *office.GetUserTagsReq) (*office.GetUserTagsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) CreateTag(ctx context.Context, req *office.CreateTagReq) (*office.CreateTagResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) DeleteTag(ctx context.Context, req *office.DeleteTagReq) (*office.DeleteTagResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) SetTag(ctx context.Context, req *office.SetTagReq) (*office.SetTagResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) SendMsg2Tag(ctx context.Context, req *office.SendMsg2TagReq) (*office.SendMsg2TagResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) GetTagSendLogs(ctx context.Context, req *office.GetTagSendLogsReq) (*office.GetTagSendLogsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) GetUserTagByID(ctx context.Context, req *office.GetUserTagByIDReq) (*office.GetUserTagByIDResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) CreateOneWorkMoment(ctx context.Context, req *office.CreateOneWorkMomentReq) (*office.CreateOneWorkMomentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) DeleteOneWorkMoment(ctx context.Context, req *office.DeleteOneWorkMomentReq) (*office.DeleteOneWorkMomentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) LikeOneWorkMoment(ctx context.Context, req *office.LikeOneWorkMomentReq) (*office.LikeOneWorkMomentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) CommentOneWorkMoment(ctx context.Context, req *office.CommentOneWorkMomentReq) (*office.CommentOneWorkMomentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) DeleteComment(ctx context.Context, req *office.DeleteCommentReq) (*office.DeleteCommentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) GetWorkMomentByID(ctx context.Context, req *office.GetWorkMomentByIDReq) (*office.GetWorkMomentByIDResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) ChangeWorkMomentPermission(ctx context.Context, req *office.ChangeWorkMomentPermissionReq) (*office.ChangeWorkMomentPermissionResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) GetUserWorkMoments(ctx context.Context, req *office.GetUserWorkMomentsReq) (*office.GetUserWorkMomentsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) GetUserFriendWorkMoments(ctx context.Context, req *office.GetUserFriendWorkMomentsReq) (*office.GetUserFriendWorkMomentsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *officeServer) SetUserWorkMomentsLevel(ctx context.Context, req *office.SetUserWorkMomentsLevelReq) (*office.SetUserWorkMomentsLevelResp, error) {
	//TODO implement me
	panic("implement me")
}
