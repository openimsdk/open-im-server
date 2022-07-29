package fcm

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/tools"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"path/filepath"
	"strconv"
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

func (f *Fcm) Push(accounts []string, alert, detailContent, operationID string, opts push.PushOpts) (string, error) {
	// accounts->registrationToken
	Tokens := make([]string, 0)
	for _, account := range accounts {
		IosfcmToken, IosErr := db.DB.GetFcmToken(account, 1)
		AndroidfcmToken, AndroidErr := db.DB.GetFcmToken(account, 2)

		if IosErr == nil {
			Tokens = append(Tokens, IosfcmToken)
		}
		if AndroidErr == nil {
			Tokens = append(Tokens, AndroidfcmToken)
		}
	}
	Success := 0
	Fail := 0
	result := tools.NewSplitter(SinglePushCountLimit, Tokens).GetSplitResult()
	Msg := new(messaging.MulticastMessage)
	Msg.Notification = &messaging.Notification{}
	Msg.Notification.Body = detailContent
	Msg.Notification.Title = alert
	ctx := context.Background()
	for _, v := range result {
		Msg.Tokens = v.Item
		//SendMulticast sends the given multicast message to all the FCM registration tokens specified.
		//The tokens array in MulticastMessage may contain up to 500 tokens.
		//SendMulticast uses the `SendAll()` function to send the given message to all the target recipients.
		//The responses list obtained from the return value corresponds to the order of the input tokens.
		//An error from SendMulticast indicates a total failure -- i.e.
		//the message could not be sent to any of the recipients.
		//Partial failures are indicated by a `BatchResponse` return value.
		response, err := f.FcmMsgCli.SendMulticast(ctx, Msg)
		if err != nil {
			Fail = Fail + len(v.Item)
			log.Info(operationID, "some token push err", err.Error(), len(v.Item))
			continue
		}
		Success = Success + response.SuccessCount
		Fail = Fail + response.FailureCount
	}
	return strconv.Itoa(Success) + " Success," + strconv.Itoa(Fail) + " Fail", nil
}
