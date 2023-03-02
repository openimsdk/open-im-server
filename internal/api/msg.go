package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/msg"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewMsg(zk *openKeeper.ZkClient) *Conversation {
	return &Conversation{zk: zk}
}

type Msg struct {
	zk *openKeeper.ZkClient
}

func (o *Msg) client() (msg.MsgClient, error) {
	conn, err := o.zk.GetConn(config.Config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		return nil, err
	}
	return msg.NewMsgClient(conn), nil
}

func (o *Msg) GetSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSeq, o.client, c)
}

func (o *Msg) SendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.SendMsg, o.client, c)
}

func (o *Msg) PullMsgBySeqList(c *gin.Context) {
	a2r.Call(msg.MsgClient.PullMsgBySeqList, o.client, c)
}

func (o *Msg) DelMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.DelMsg, o.client, c)
}

func (o *Msg) DelSuperGroupMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.DelSuperGroupMsg, o.client, c)
}

func (o *Msg) ClearMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ClearMsg, o.client, c)
}

func (o *Msg) SetMsgMinSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetMsgMinSeq, o.client, c)
}

func (o *Msg) SetMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetMessageReactionExtensions, o.client, c)
}

func (o *Msg) GetMessageListReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMessageListReactionExtensions, o.client, c)
}

func (o *Msg) AddMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.AddMessageReactionExtensions, o.client, c)
}

func (o *Msg) DeleteMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMessageReactionExtensions, o.client, c)
}

func (o *Msg) ManagementSendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ManagementSendMsg, o.client, c)
}

func (o *Msg) ManagementBatchSendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ManagementBatchSendMsg, o.client, c)
}

func (o *Msg) CheckMsgIsSendSuccess(c *gin.Context) {
	a2r.Call(msg.MsgClient.CheckMsgIsSendSuccess, o.client, c)
}
