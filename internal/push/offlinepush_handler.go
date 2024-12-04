package push

import (
	"context"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/protocol/constant"
	pbpush "github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/utils/jsonutil"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
)

type OfflinePushConsumerHandler struct {
	OfflinePushConsumerGroup *kafka.MConsumerGroup
	offlinePusher            offlinepush.OfflinePusher
}

func NewOfflinePushConsumerHandler(config *Config, offlinePusher offlinepush.OfflinePusher) (*OfflinePushConsumerHandler, error) {
	var offlinePushConsumerHandler OfflinePushConsumerHandler
	var err error
	offlinePushConsumerHandler.offlinePusher = offlinePusher
	offlinePushConsumerHandler.OfflinePushConsumerGroup, err = kafka.NewMConsumerGroup(config.KafkaConfig.Build(), config.KafkaConfig.ToOfflineGroupID,
		[]string{config.KafkaConfig.ToOfflinePushTopic}, true)
	if err != nil {
		return nil, err
	}
	return &offlinePushConsumerHandler, nil
}

func (*OfflinePushConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*OfflinePushConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (o *OfflinePushConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := o.OfflinePushConsumerGroup.GetContextFromMsg(msg)
		o.handleMsg2OfflinePush(ctx, msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (o *OfflinePushConsumerHandler) handleMsg2OfflinePush(ctx context.Context, msg []byte) {
	offlinePushMsg := pbpush.PushMsgReq{}
	if err := proto.Unmarshal(msg, &offlinePushMsg); err != nil {
		log.ZError(ctx, "offline push Unmarshal msg err", err, "msg", string(msg))
		return
	}
	if offlinePushMsg.MsgData == nil || offlinePushMsg.UserIDs == nil {
		log.ZError(ctx, "offline push msg is empty", errs.New("offlinePushMsg is empty"), "userIDs", offlinePushMsg.UserIDs, "msg", offlinePushMsg.MsgData)
		return
	}
	if offlinePushMsg.MsgData.Status == constant.MsgStatusSending {
		offlinePushMsg.MsgData.Status = constant.MsgStatusSendSuccess
	}
	log.ZInfo(ctx, "receive to OfflinePush MQ", "userIDs", offlinePushMsg.UserIDs, "msg", offlinePushMsg.MsgData)

	err := o.offlinePushMsg(ctx, offlinePushMsg.MsgData, offlinePushMsg.UserIDs)
	if err != nil {
		log.ZWarn(ctx, "offline push failed", err, "msg", offlinePushMsg.String())
	}
}

func (o *OfflinePushConsumerHandler) getOfflinePushInfos(msg *sdkws.MsgData) (title, content string, opts *options.Opts, err error) {
	type AtTextElem struct {
		Text       string   `json:"text,omitempty"`
		AtUserList []string `json:"atUserList,omitempty"`
		IsAtSelf   bool     `json:"isAtSelf"`
	}

	opts = &options.Opts{Signal: &options.Signal{ClientMsgID: msg.ClientMsgID}}
	if msg.OfflinePushInfo != nil {
		opts.IOSBadgeCount = msg.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = msg.OfflinePushInfo.IOSPushSound
		opts.Ex = msg.OfflinePushInfo.Ex
	}

	if msg.OfflinePushInfo != nil {
		title = msg.OfflinePushInfo.Title
		content = msg.OfflinePushInfo.Desc
	}
	if title == "" {
		switch msg.ContentType {
		case constant.Text:
			fallthrough
		case constant.Picture:
			fallthrough
		case constant.Voice:
			fallthrough
		case constant.Video:
			fallthrough
		case constant.File:
			title = constant.ContentType2PushContent[int64(msg.ContentType)]
		case constant.AtText:
			ac := AtTextElem{}
			_ = jsonutil.JsonStringToStruct(string(msg.Content), &ac)
		case constant.SignalingNotification:
			title = constant.ContentType2PushContent[constant.SignalMsg]
		default:
			title = constant.ContentType2PushContent[constant.Common]
		}
	}
	if content == "" {
		content = title
	}
	return
}

func (o *OfflinePushConsumerHandler) offlinePushMsg(ctx context.Context, msg *sdkws.MsgData, offlinePushUserIDs []string) error {
	title, content, opts, err := o.getOfflinePushInfos(msg)
	if err != nil {
		return err
	}
	err = o.offlinePusher.Push(ctx, offlinePushUserIDs, title, content, opts)
	if err != nil {
		prommetrics.MsgOfflinePushFailedCounter.Inc()
		return err
	}
	return nil
}
