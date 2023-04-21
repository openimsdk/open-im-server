package notification2

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte //  sdkws.TipsComm
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	SenderNickname string
	SenderFaceURL  string
}
