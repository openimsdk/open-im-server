/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 15:23).
 */
package manage

import "github.com/gin-gonic/gin"

//
//var validate *validator.Validate
//
//
//func newUserSendMsgReq(params *paramsManagementSendMsg) *pbChat.SendMsgReq {
//	var newContent string
//	switch params.ContentType {
//	case constant.Text:
//		newContent = params.Content["text"].(string)
//	case constant.Picture:
//		fallthrough
//	case constant.Custom:
//		fallthrough
//	case constant.Voice:
//		fallthrough
//	case constant.File:
//		newContent = utils.StructToJsonString(params.Content)
//	default:
//	}
//	options := make(map[string]bool, 2)
//	if params.IsOnlineOnly {
//		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
//		utils.SetSwitchFromOptions(options, constant.IsHistory, false)
//		utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
//	}
//	pbData := pbChat.SendMsgReq{
//		OperationID: params.OperationID,
//		MsgData: &open_im_sdk.MsgData{
//			SendID:           params.SendID,
//			RecvID:           params.RecvID,
//			GroupID:          params.GroupID,
//			ClientMsgID:      utils.GetMsgID(params.SendID),
//			SenderPlatformID: params.SenderPlatformID,
//			SenderNickName:   params.SenderNickName,
//			SenderFaceURL:    params.SenderFaceURL,
//			SessionType:      params.SessionType,
//			MsgFrom:          constant.SysMsgType,
//			ContentType:      params.ContentType,
//			Content:          []byte(newContent),
//			ForceList:        params.ForceList,
//			CreateTime:       utils.GetCurrentTimestampByNano(),
//			Options:          options,
//			OfflinePushInfo:  params.OfflinePushInfo,
//		},
//	}
//	return &pbData
//}
//func init() {
//	validate = validator.New()
//}
func ManagementSendMsg(c *gin.Context) {

}

//func ManagementSendMsg(c *gin.Context) {
//	var data interface{}
//	params := paramsManagementSendMsg{}
//	if err := c.BindJSON(&params); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
//		log.ErrorByKv("json unmarshal err", c.PostForm("operationID"), "err", err.Error(), "content", c.PostForm("content"))
//		return
//	}
//	switch params.ContentType {
//	case constant.Text:
//		data = TextElem{}
//	case constant.Picture:
//		data = PictureElem{}
//	case constant.Voice:
//		data = SoundElem{}
//	case constant.Video:
//		data = VideoElem{}
//	case constant.File:
//		data = FileElem{}
//	//case constant.AtText:
//	//	data = AtElem{}
//	//case constant.Merger:
//	//	data =
//	//case constant.Card:
//	//case constant.Location:
//	case constant.Custom:
//		data = CustomElem{}
//	//case constant.Revoke:
//	//case constant.HasReadReceipt:
//	//case constant.Typing:
//	//case constant.Quote:
//	default:
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 404, "errMsg": "contentType err"})
//		log.ErrorByKv("contentType err", c.PostForm("operationID"), "content", c.PostForm("content"))
//		return
//	}
//	if err := mapstructure.WeakDecode(params.Content, &data); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": err.Error()})
//		log.ErrorByKv("content to Data struct  err", "", "err", err.Error())
//		return
//	} else if err := validate.Struct(data); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 403, "errMsg": err.Error()})
//		log.ErrorByKv("data args validate  err", "", "err", err.Error())
//		return
//	}
//
//	token := c.Request.Header.Get("token")
//	claims, err := token_verify.ParseToken(token)
//	if err != nil {
//		log.NewError(params.OperationID, "parse token failed", err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parse token failed", "sendTime": 0, "MsgID": ""})
//	}
//	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "not authorized", "sendTime": 0, "MsgID": ""})
//		return
//
//	}
//	switch params.SessionType {
//	case constant.SingleChatType:
//		if len(params.RecvID) == 0 {
//			log.NewError(params.OperationID, "recvID is a null string")
//			c.JSON(http.StatusBadRequest, gin.H{"errCode": 405, "errMsg": "recvID is a null string", "sendTime": 0, "MsgID": ""})
//		}
//	case constant.GroupChatType:
//		if len(params.GroupID) == 0 {
//			log.NewError(params.OperationID, "groupID is a null string")
//			c.JSON(http.StatusBadRequest, gin.H{"errCode": 405, "errMsg": "groupID is a null string", "sendTime": 0, "MsgID": ""})
//		}
//
//	}
//	log.InfoByKv("Ws call success to ManagementSendMsgReq", params.OperationID, "Parameters", params)
//
//	pbData := newUserSendMsgReq(&params)
//	log.Info("", "", "api ManagementSendMsg call start..., [data: %s]", pbData.String())
//
//	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
//	client := pbChat.NewChatClient(etcdConn)
//
//	log.Info("", "", "api ManagementSendMsg call, api call rpc...")
//
//	reply, err := client.SendMsg(context.Background(), pbData)
//	if err != nil {
//		log.NewError(params.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call UserSendMsg  rpc server failed"})
//		return
//	}
//	log.Info("", "", "api ManagementSendMsg call end..., [data: %s] [reply: %s]", pbData.String(), reply.String())
//
//	c.JSON(http.StatusOK, gin.H{
//		"errCode":  reply.ErrCode,
//		"errMsg":   reply.ErrMsg,
//		"sendTime": reply.SendTime,
//		"msgID":    reply.ClientMsgID,
//	})
//
//}

//type MergeElem struct {
//	Title        string       `json:"title"`
//	AbstractList []string     `json:"abstractList"`
//	MultiMessage []*MsgStruct `json:"multiMessage"`
//}

//type QuoteElem struct {
//	Text         string     `json:"text"`
//	QuoteMessage *MsgStruct `json:"quoteMessage"`
//}
