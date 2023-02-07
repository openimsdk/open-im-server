package group

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/tracelog"
	pbConversation "Open_IM/pkg/proto/conversation"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"math/big"
	"strconv"
	"time"
)

func GetPublicUserInfoOne(ctx context.Context, userID string) (*sdk_ws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdk_ws.UserInfo, error) {
	return nil, errors.New("todo")
}

func GetUserInfoMap(ctx context.Context, userIDs []string) (map[string]*sdk_ws.UserInfo, error) {
	users, err := GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdk_ws.UserInfo) string {
		return e.UserID
	}), nil
}

func GetPublicUserInfo(ctx context.Context, userIDs []string) ([]*sdk_ws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func GetPublicUserInfoMap(ctx context.Context, userIDs []string) (map[string]*sdk_ws.PublicUserInfo, error) {
	users, err := GetPublicUserInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdk_ws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

func GroupNotification(ctx context.Context, groupID string) {
	var conversationReq pbConversation.ModifyConversationFieldReq
	conversation := pbConversation.Conversation{
		OwnerUserID:      tracelog.GetOpUserID(ctx),
		ConversationID:   utils.GetConversationIDBySessionType(groupID, constant.GroupChatType),
		ConversationType: constant.GroupChatType,
		GroupID:          groupID,
	}
	conversationReq.Conversation = &conversation
	conversationReq.OperationID = tracelog.GetOperationID(ctx)
	conversationReq.FieldType = constant.FieldGroupAtType
	conversation.GroupAtType = constant.GroupNotification
	conversationReq.UserIDList = cacheResp.UserIDList

	_, err = pbConversation.NewConversationClient(s.etcdConn.GetConn("", config.Config.RpcRegisterName.OpenImConversationName)).ModifyConversationField(ctx, &conversationReq)
	tracelog.SetCtxInfo(ctx, "ModifyConversationField", err, "req", &conversationReq, "resp", conversationReply)
}

func genGroupID(ctx context.Context, groupID string) string {
	if groupID != "" {
		return groupID
	}
	groupID = utils.Md5(tracelog.GetOperationID(ctx) + strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(groupID[0:8], 16)
	groupID = bi.String()
	return groupID
}
