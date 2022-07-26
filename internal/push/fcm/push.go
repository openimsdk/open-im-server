package push

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"context"
	"log"
	"path/filepath"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Fcm struct {
}

var (
	FcmClient *Fcm
	FcmMsgCli *messaging.Client
)

func init() {
	//FcmClient = newFcmClient()
}

func newFcmClient() *Fcm {
	opt := option.WithCredentialsFile(filepath.Join(config.Root, "config", config.Config.Push.Fcm.ServiceAccount))
	fcmApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("error initializing app: %v\n", err)
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
	FcmMsgCli, err = fcmApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return nil
	}
	log.Println(FcmMsgCli)
	return &Fcm{}
}

func (f *Fcm) Push(accounts []string, alert, detailContent, operationID string, opts push.PushOpts) (string, error) {
	//需要一个客户端的Token
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
	tokenlen := len(Tokens)
	// 500组为一个推送，我们用400好了
	limit := 400
	pages := int((tokenlen-1)/limit + 1)
	Success := 0
	Fail := 0
	for i := 0; i < pages; i++ {
		Msg := new(messaging.MulticastMessage)
		Msg.Notification.Body = detailContent
		Msg.Notification.Title = alert
		ctx := context.Background()
		max := (i+1)*limit - 1
		if max >= tokenlen {
			max = tokenlen - 1
		}
		Msg.Tokens = Tokens[i*limit : max]
		//SendMulticast sends the given multicast message to all the FCM registration tokens specified.
		//The tokens array in MulticastMessage may contain up to 500 tokens.
		//SendMulticast uses the `SendAll()` function to send the given message to all the target recipients.
		//The responses list obtained from the return value corresponds to the order of the input tokens.
		//An error from SendMulticast indicates a total failure -- i.e.
		//the message could not be sent to any of the recipients.
		//Partial failures are indicated by a `BatchResponse` return value.
		response, err := FcmMsgCli.SendMulticast(ctx, Msg)
		if err != nil {
			log.Fatalln(err)
		}
		Success = Success + response.SuccessCount
		Fail = Fail + response.FailureCount
	}
	return strconv.Itoa(Success) + " Success," + strconv.Itoa(Fail) + " Fail", nil
}
