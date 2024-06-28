package msggateway

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"math/rand"
	"strconv"
	"time"
)

func (ws *WsServer) ChangeOnlineStatus(concurrent int) {
	if concurrent < 1 {
		concurrent = 1
	}
	scanTicker := time.NewTicker(time.Minute * 3)

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

	opIdCtx := mcontext.SetOperationID(context.Background(), "r"+strconv.FormatUint(rNum, 10))
	doRequest := func(req *pbuser.SetUserOnlineStatusReq) {
		ctx, cancel := context.WithTimeout(opIdCtx, time.Second*5)
		defer cancel()
		if _, err := ws.userClient.Client.SetUserOnlineStatus(ctx, req); err != nil {
			log.ZError(ctx, "update user online status", err)
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
		case now := <-scanTicker.C:
			pushUserState(ws.clients.GetAllUserStatus(now.Add(-cachekey.OnlineExpire/3), now)...)
		case state := <-ws.clients.UserState():
			log.ZDebug(context.Background(), "user online change", "userID", state.UserID, "online", state.Online, "offline", state.Offline)
			pushUserState(state)
		}
	}
}
