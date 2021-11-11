package req

import (
	tpns "Open_IM/internal/push/sdk/tpns-server-sdk-go/go"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var PushURL = "https://api.tpns.tencent.com/v3/push/app"

//var PushURL = "https://test.api.tpns.tencent.com/v3/push/app"

func URL(url string) {
	PushURL = url
}

type ReqOpt func(*tpns.Request)

func NewPush(req *tpns.Request, opts ...ReqOpt) (*http.Request, string, error) {
	return NewPushReq(req, opts...)
}

func NewUploadFileRequest(host string, file string) (*http.Request, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(fp.Name()))
	if err != nil {
		return nil, err
	}

	io.Copy(part, fp)
	writer.Close()
	url := host + "/v3/push/package/upload"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func NewSingleAccountPush(
	message tpns.Message,
	account string,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdAccountList,
		AccountList:  []string{account},
		Message:      message,
	}
	return NewPushReq(req, opts...)
}

func NewListAccountPush(
	accounts []string, message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdAccountList,
		AccountList:  accounts,
		Message:      message,
		Environment:  tpns.EnvDev,
	}
	return NewPushReq(req, opts...)
}

func NewTokenPush(
	tokens []string, message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdTokenList,
		TokenList:    tokens,
		Message:      message,
		Environment:  tpns.EnvProd,
	}
	//fmt.Printf("reqBody :%v", common.ToJson(req))
	//fmt.Println()
	return NewPushReq(req, opts...)
}

func NewTagsPush(
	tagList []tpns.TagRule, message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdTag,
		Tag:          tagList,
		Message:      message,
	}
	//fmt.Printf("reqBody :%v", common.ToJson(req))
	//fmt.Println()
	return NewPushReq(req, opts...)
}

func NewAllPush(
	message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdAll,
		Message:      message,
	}
	return NewPushReq(req, opts...)
}

func NewAccountPackagePush(
	message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdPackageAccount,
		Message:      message,
	}
	return NewPushReq(req, opts...)
}

func NewTokenPackagePush(
	message tpns.Message,
	opts ...ReqOpt,
) (*http.Request, string, error) {
	req := &tpns.Request{
		MessageType:  tpns.MsgTypeNotify,
		AudienceType: tpns.AdPackageToken,
		Message:      message,
	}
	return NewPushReq(req, opts...)
}

func NewPushReq(req *tpns.Request, opts ...ReqOpt) (request *http.Request, reqBody string, err error) {
	for _, opt := range opts {
		opt(req)
	}
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}
	reqBody = string(bodyBytes)
	//fmt.Printf("NewPushReq req:%v", reqBody)
	request, err = http.NewRequest("POST", PushURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", err
	}
	request.Header.Add("Content-Type", "application/json")
	return
}

func EnvProd() ReqOpt {
	return func(r *tpns.Request) {
		r.Environment = tpns.EnvProd
	}
}

func EnvDev() ReqOpt {
	return func(r *tpns.Request) {
		r.Environment = tpns.EnvDev
	}
}

func Title(t string) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Title = t
		if r.Message.IOS != nil {
			if r.Message.IOS.Aps != nil {
				r.Message.IOS.Aps.Alert["title"] = t
			} else {
				r.Message.IOS.Aps = &tpns.Aps{
					Alert: map[string]string{"title": t},
				}
			}
		} else {
			r.Message.IOS = &tpns.IOSParams{
				Aps: &tpns.Aps{
					Alert: map[string]string{"title": t},
				},
			}
		}
	}
}

func Content(c string) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Content = c
		if r.Message.IOS != nil {
			if r.Message.IOS.Aps != nil {
				r.Message.IOS.Aps.Alert["body"] = c
			} else {
				r.Message.IOS.Aps = &tpns.Aps{
					Alert: map[string]string{"body": c},
				}
			}
		} else {
			r.Message.IOS = &tpns.IOSParams{
				Aps: &tpns.Aps{
					Alert: map[string]string{"body": c},
				},
			}
		}
	}
}

func Ring(ring *int) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.Ring = ring
	}
}

func RingRaw(rr string) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.RingRaw = rr
	}
}

func Vibrate(v *int) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.Vibrate = v
	}
}

func Lights(l *int) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.Lights = l
	}
}

func Clearable(c *int) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.Clearable = c
	}
}

func IconType(it *int) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.IconType = it
	}
}

func IconRes(ir string) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.IconRes = ir
	}
}

func AndroidCustomContent(ct string) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.Android.CustomContent = ct
	}
}

func Aps(aps *tpns.Aps) ReqOpt {
	return func(r *tpns.Request) {
		r.Message.IOS.Aps = aps
	}
}

func AudienceType(at tpns.AudienceType) ReqOpt {
	return func(r *tpns.Request) {
		r.AudienceType = at
	}
}

func Message(m tpns.Message) ReqOpt {
	return func(r *tpns.Request) {
		r.Message = m
	}
}

func TokenList(tl []string) ReqOpt {
	return func(r *tpns.Request) {
		r.TokenList = tl
	}
}

func TokenListAdd(t string) ReqOpt {
	return func(r *tpns.Request) {
		if r.TokenList != nil {
			r.TokenList = append(r.TokenList, t)
		} else {
			r.TokenList = []string{t}
		}
	}
}

func AccountList(al []string) ReqOpt {
	return func(r *tpns.Request) {
		r.AccountList = al
	}
}

//ChannelDistributeRules
func AddChannelRules(ChannelRules []*tpns.ChannelDistributeRule) ReqOpt {
	return func(r *tpns.Request) {
		r.ChannelRules = ChannelRules
	}
}

//ChannelDistributeRules
func AddLoopParam(loopParam *tpns.PushLoopParam) ReqOpt {
	return func(r *tpns.Request) {
		r.LoopParam = loopParam
	}
}

func AccountListAdd(a string) ReqOpt {
	return func(r *tpns.Request) {
		if r.AccountList != nil {
			r.AccountList = append(r.AccountList, a)
		} else {
			r.AccountList = []string{a}
		}
	}
}

func MessageType(t tpns.MessageType) ReqOpt {
	return func(r *tpns.Request) {
		r.MessageType = t
	}
}

func AddMultiPkg(multipPkg bool) ReqOpt {
	return func(r *tpns.Request) {
		r.MultiPkg = multipPkg
	}
}

func AddForceCollapse(forceCollapse bool) ReqOpt {
	return func(r *tpns.Request) {
		r.ForceCollapse = forceCollapse
	}
}

func AddTPNSOnlinePushType(onlinePushType int) ReqOpt {
	return func(r *tpns.Request) {
		r.TPNSOnlinePushType = onlinePushType
	}
}

func AddCollapseId(collapseId int) ReqOpt {
	return func(r *tpns.Request) {
		r.CollapseId = collapseId
	}
}

func AddPushSpeed(pushSpeed int) ReqOpt {
	return func(r *tpns.Request) {
		r.PushSpeed = pushSpeed
	}
}

func AddAccountPushType(accountPushType int) ReqOpt {
	return func(r *tpns.Request) {
		r.AccountPushType = accountPushType
	}
}

func AddPlanId(planId string) ReqOpt {
	return func(r *tpns.Request) {
		r.PlanId = planId
	}
}

func AddSendTime(sendTime string) ReqOpt {
	return func(r *tpns.Request) {
		r.SendTime = sendTime
	}
}

func AddExpireTime(expireTime int) ReqOpt {
	return func(r *tpns.Request) {
		r.ExpireTime = expireTime
	}
}

func AddUploadId(UploadId int) ReqOpt {
	return func(r *tpns.Request) {
		r.UploadId = UploadId
	}
}

func AddEnvironment(Environment tpns.CommonRspEnv) ReqOpt {
	return func(r *tpns.Request) {
		r.Environment = Environment
	}
}
