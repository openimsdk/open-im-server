package fcm

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"context"
	go_redis "github.com/go-redis/redis/v8"
	"path/filepath"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

const SinglePushCountLimit = 400

type Fcm struct {
	FcmMsgCli *messaging.Client
}

func NewFcm() *Fcm {
	return newFcmClient()
}
func newFcmClient() *Fcm {
	opt := option.WithCredentialsFile(filepath.Join(config.Root, "config", config.Config.Push.Fcm.ServiceAccount))
	fcmApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Debug("", "error initializing app: ", err.Error())
		return nil
	}
	//授权
	// fcmClient, err := fcmApp.Auth(context.Background())
	// if err != nil {
	// 	log.Println("error getting Auth client: %v\n", err)
	// 	return
	// }
	// log.Printf("%#v\r\n", fcmClient)
	ctx := context.Background()
	fcmMsgClient, err := fcmApp.Messaging(ctx)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return &Fcm{FcmMsgCli: fcmMsgClient}
}

func (f *Fcm) Push(accounts []string, title, detailContent, operationID string, opts push.PushOpts) (string, error) {
	// accounts->registrationToken
	allTokens := make(map[string][]string, 0)
	for _, account := range accounts {
		var personTokens []string
		for _, v := range push.PushTerminal {
			Token, err := db.DB.GetFcmToken(account, v)
			if err == nil {
				personTokens = append(personTokens, Token)
			}
		}
		allTokens[account] = personTokens
	}
	Success := 0
	Fail := 0
	notification := &messaging.Notification{}
	notification.Body = detailContent
	notification.Title = title
	var messages []*messaging.Message
	ctx := context.Background()
	for uid, personTokens := range allTokens {
		apns := &messaging.APNSConfig{Payload: &messaging.APNSPayload{Aps: &messaging.Aps{Sound: opts.IOSPushSound}}}
		messageCount := len(messages)
		if messageCount >= SinglePushCountLimit {
			response, err := f.FcmMsgCli.SendAll(ctx, messages)
			if err != nil {
				Fail = Fail + messageCount
				log.Info(operationID, "some token push err", err.Error(), messageCount)
			} else {
				Success = Success + response.SuccessCount
				Fail = Fail + response.FailureCount
			}
			messages = messages[0:0]
		}
		if opts.IOSBadgeCount {
			unreadCountSum, err := db.DB.IncrUserBadgeUnreadCountSum(uid)
			if err == nil {
				apns.Payload.Aps.Badge = &unreadCountSum
			} else {
				log.Error(operationID, "IncrUserBadgeUnreadCountSum redis err", err.Error(), uid)
				Fail++
				continue
			}
		} else {
			unreadCountSum, err := db.DB.GetUserBadgeUnreadCountSum(uid)
			if err == nil && unreadCountSum != 0 {
				apns.Payload.Aps.Badge = &unreadCountSum
			} else if err == go_redis.Nil || unreadCountSum == 0 {
				zero := 1
				apns.Payload.Aps.Badge = &zero
			} else {
				log.Error(operationID, "GetUserBadgeUnreadCountSum redis err", err.Error(), uid)
				Fail++
				continue
			}
		}
		for _, token := range personTokens {
			temp := &messaging.Message{
				Data:         map[string]string{"ex": opts.Data},
				Token:        token,
				Notification: notification,
				APNS:         apns,
			}
			messages = append(messages, temp)
		}

	}
	messageCount := len(messages)
	if messageCount > 0 {
		response, err := f.FcmMsgCli.SendAll(ctx, messages)
		if err != nil {
			Fail = Fail + messageCount
			log.Info(operationID, "some token push err", err.Error(), messageCount)
		} else {
			Success = Success + response.SuccessCount
			Fail = Fail + response.FailureCount
		}
	}
	return strconv.Itoa(Success) + " Success," + strconv.Itoa(Fail) + " Fail", nil
}
