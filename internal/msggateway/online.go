package msggateway

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func (ws *WsServer) ChangeOnlineStatus(concurrent int) {
	if concurrent < 1 {
		concurrent = 1
	}
	const renewalTime = cachekey.OnlineExpire / 3
	//const renewalTime = time.Second * 10
	renewalTicker := time.NewTicker(renewalTime)

	requestChs := make([]chan *pbuser.SetUserOnlineStatusReq, concurrent)
	changeStatus := make([][]UserState, concurrent)

	for i := 0; i < concurrent; i++ {
		requestChs[i] = make(chan *pbuser.SetUserOnlineStatusReq, 64)
		changeStatus[i] = make([]UserState, 0, 100)
	}

	mergeTicker := time.NewTicker(time.Second)

	local2pb := func(u UserState) *pbuser.UserOnlineStatus {
		return &pbuser.UserOnlineStatus{
			UserID:  u.UserID,
			Online:  u.Online,
			Offline: u.Offline,
		}
	}

	rNum := rand.Uint64()
	pushUserState := func(us ...UserState) {
		for _, u := range us {
			sum := md5.Sum([]byte(u.UserID))
			i := (binary.BigEndian.Uint64(sum[:]) + rNum) % uint64(concurrent)
			changeStatus[i] = append(changeStatus[i], u)
			status := changeStatus[i]
			if len(status) == cap(status) {
				req := &pbuser.SetUserOnlineStatusReq{
					Status: datautil.Slice(status, local2pb),
				}
				changeStatus[i] = status[:0]
				select {
				case requestChs[i] <- req:
				default:
					log.ZError(context.Background(), "user online processing is too slow", nil)
				}
			}
		}
	}

	pushAllUserState := func() {
		for i, status := range changeStatus {
			if len(status) == 0 {
				continue
			}
			req := &pbuser.SetUserOnlineStatusReq{
				Status: datautil.Slice(status, local2pb),
			}
			changeStatus[i] = status[:0]
			select {
			case requestChs[i] <- req:
			default:
				log.ZError(context.Background(), "user online processing is too slow", nil)
			}
		}
	}

	var count atomic.Int64
	operationIDPrefix := fmt.Sprintf("p_%d_", os.Getpid())
	doRequest := func(req *pbuser.SetUserOnlineStatusReq) {
		opIdCtx := mcontext.SetOperationID(context.Background(), operationIDPrefix+strconv.FormatInt(count.Add(1), 10))
		ctx, cancel := context.WithTimeout(opIdCtx, time.Second*5)
		defer cancel()
		if _, err := ws.userClient.Client.SetUserOnlineStatus(ctx, req); err != nil {
			log.ZError(ctx, "update user online status", err)
		}
		for _, ss := range req.Status {
			for _, online := range ss.Online {
				client, _, _ := ws.clients.Get(ss.UserID, int(online))
				back := false
				if len(client) > 0 {
					back = client[0].IsBackground
				}
				ws.webhookAfterUserOnline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOnline, ss.UserID, int(online), back, ss.ConnID)
			}
			for _, offline := range ss.Offline {
				ws.webhookAfterUserOffline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOffline, ss.UserID, int(offline), ss.ConnID)
			}
		}
	}

	for i := 0; i < concurrent; i++ {
		go func(ch <-chan *pbuser.SetUserOnlineStatusReq) {
			for req := range ch {
				doRequest(req)
			}
		}(requestChs[i])
	}

	for {
		select {
		case <-mergeTicker.C:
			pushAllUserState()
		case now := <-renewalTicker.C:
			deadline := now.Add(-cachekey.OnlineExpire / 3)
			users := ws.clients.GetAllUserStatus(deadline, now)
			log.ZDebug(context.Background(), "renewal ticker", "deadline", deadline, "nowtime", now, "num", len(users), "users", users)
			pushUserState(users...)
		case state := <-ws.clients.UserState():
			log.ZDebug(context.Background(), "OnlineCache user online change", "userID", state.UserID, "online", state.Online, "offline", state.Offline)
			pushUserState(state)
		}
	}
}
