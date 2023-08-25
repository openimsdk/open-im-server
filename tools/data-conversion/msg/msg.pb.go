// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v4.22.0
// source: msg/msg.proto

package msg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MsgDataToMQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token       string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	OperationID string   `protobuf:"bytes,2,opt,name=operationID,proto3" json:"operationID"`
	MsgData     *MsgData `protobuf:"bytes,3,opt,name=msgData,proto3" json:"msgData"`
}

func (x *MsgDataToMQ) Reset() {
	*x = MsgDataToMQ{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgDataToMQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgDataToMQ) ProtoMessage() {}

func (x *MsgDataToMQ) ProtoReflect() protoreflect.Message {
	mi := &file_msg_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgDataToMQ.ProtoReflect.Descriptor instead.
func (*MsgDataToMQ) Descriptor() ([]byte, []int) {
	return file_msg_msg_proto_rawDescGZIP(), []int{0}
}

func (x *MsgDataToMQ) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *MsgDataToMQ) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *MsgDataToMQ) GetMsgData() *MsgData {
	if x != nil {
		return x.MsgData
	}
	return nil
}

type MsgData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SendID               string           `protobuf:"bytes,1,opt,name=sendID,proto3" json:"sendID"`
	RecvID               string           `protobuf:"bytes,2,opt,name=recvID,proto3" json:"recvID"`
	GroupID              string           `protobuf:"bytes,3,opt,name=groupID,proto3" json:"groupID"`
	ClientMsgID          string           `protobuf:"bytes,4,opt,name=clientMsgID,proto3" json:"clientMsgID"`
	ServerMsgID          string           `protobuf:"bytes,5,opt,name=serverMsgID,proto3" json:"serverMsgID"`
	SenderPlatformID     int32            `protobuf:"varint,6,opt,name=senderPlatformID,proto3" json:"senderPlatformID"`
	SenderNickname       string           `protobuf:"bytes,7,opt,name=senderNickname,proto3" json:"senderNickname"`
	SenderFaceURL        string           `protobuf:"bytes,8,opt,name=senderFaceURL,proto3" json:"senderFaceURL"`
	SessionType          int32            `protobuf:"varint,9,opt,name=sessionType,proto3" json:"sessionType"`
	MsgFrom              int32            `protobuf:"varint,10,opt,name=msgFrom,proto3" json:"msgFrom"`
	ContentType          int32            `protobuf:"varint,11,opt,name=contentType,proto3" json:"contentType"`
	Content              []byte           `protobuf:"bytes,12,opt,name=content,proto3" json:"content"`
	Seq                  uint32           `protobuf:"varint,14,opt,name=seq,proto3" json:"seq"`
	SendTime             int64            `protobuf:"varint,15,opt,name=sendTime,proto3" json:"sendTime"`
	CreateTime           int64            `protobuf:"varint,16,opt,name=createTime,proto3" json:"createTime"`
	Status               int32            `protobuf:"varint,17,opt,name=status,proto3" json:"status"`
	Options              map[string]bool  `protobuf:"bytes,18,rep,name=options,proto3" json:"options" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	OfflinePushInfo      *OfflinePushInfo `protobuf:"bytes,19,opt,name=offlinePushInfo,proto3" json:"offlinePushInfo"`
	AtUserIDList         []string         `protobuf:"bytes,20,rep,name=atUserIDList,proto3" json:"atUserIDList"`
	MsgDataList          []byte           `protobuf:"bytes,21,opt,name=msgDataList,proto3" json:"msgDataList"`
	AttachedInfo         string           `protobuf:"bytes,22,opt,name=attachedInfo,proto3" json:"attachedInfo"`
	Ex                   string           `protobuf:"bytes,23,opt,name=ex,proto3" json:"ex"`
	IsReact              bool             `protobuf:"varint,40,opt,name=isReact,proto3" json:"isReact"`
	IsExternalExtensions bool             `protobuf:"varint,41,opt,name=isExternalExtensions,proto3" json:"isExternalExtensions"`
	MsgFirstModifyTime   int64            `protobuf:"varint,42,opt,name=msgFirstModifyTime,proto3" json:"msgFirstModifyTime"`
}

func (x *MsgData) Reset() {
	*x = MsgData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgData) ProtoMessage() {}

func (x *MsgData) ProtoReflect() protoreflect.Message {
	mi := &file_msg_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgData.ProtoReflect.Descriptor instead.
func (*MsgData) Descriptor() ([]byte, []int) {
	return file_msg_msg_proto_rawDescGZIP(), []int{1}
}

func (x *MsgData) GetSendID() string {
	if x != nil {
		return x.SendID
	}
	return ""
}

func (x *MsgData) GetRecvID() string {
	if x != nil {
		return x.RecvID
	}
	return ""
}

func (x *MsgData) GetGroupID() string {
	if x != nil {
		return x.GroupID
	}
	return ""
}

func (x *MsgData) GetClientMsgID() string {
	if x != nil {
		return x.ClientMsgID
	}
	return ""
}

func (x *MsgData) GetServerMsgID() string {
	if x != nil {
		return x.ServerMsgID
	}
	return ""
}

func (x *MsgData) GetSenderPlatformID() int32 {
	if x != nil {
		return x.SenderPlatformID
	}
	return 0
}

func (x *MsgData) GetSenderNickname() string {
	if x != nil {
		return x.SenderNickname
	}
	return ""
}

func (x *MsgData) GetSenderFaceURL() string {
	if x != nil {
		return x.SenderFaceURL
	}
	return ""
}

func (x *MsgData) GetSessionType() int32 {
	if x != nil {
		return x.SessionType
	}
	return 0
}

func (x *MsgData) GetMsgFrom() int32 {
	if x != nil {
		return x.MsgFrom
	}
	return 0
}

func (x *MsgData) GetContentType() int32 {
	if x != nil {
		return x.ContentType
	}
	return 0
}

func (x *MsgData) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *MsgData) GetSeq() uint32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *MsgData) GetSendTime() int64 {
	if x != nil {
		return x.SendTime
	}
	return 0
}

func (x *MsgData) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *MsgData) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *MsgData) GetOptions() map[string]bool {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *MsgData) GetOfflinePushInfo() *OfflinePushInfo {
	if x != nil {
		return x.OfflinePushInfo
	}
	return nil
}

func (x *MsgData) GetAtUserIDList() []string {
	if x != nil {
		return x.AtUserIDList
	}
	return nil
}

func (x *MsgData) GetMsgDataList() []byte {
	if x != nil {
		return x.MsgDataList
	}
	return nil
}

func (x *MsgData) GetAttachedInfo() string {
	if x != nil {
		return x.AttachedInfo
	}
	return ""
}

func (x *MsgData) GetEx() string {
	if x != nil {
		return x.Ex
	}
	return ""
}

func (x *MsgData) GetIsReact() bool {
	if x != nil {
		return x.IsReact
	}
	return false
}

func (x *MsgData) GetIsExternalExtensions() bool {
	if x != nil {
		return x.IsExternalExtensions
	}
	return false
}

func (x *MsgData) GetMsgFirstModifyTime() int64 {
	if x != nil {
		return x.MsgFirstModifyTime
	}
	return 0
}

type OfflinePushInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title         string `protobuf:"bytes,1,opt,name=title,proto3" json:"title"`
	Desc          string `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc"`
	Ex            string `protobuf:"bytes,3,opt,name=ex,proto3" json:"ex"`
	IOSPushSound  string `protobuf:"bytes,4,opt,name=iOSPushSound,proto3" json:"iOSPushSound"`
	IOSBadgeCount bool   `protobuf:"varint,5,opt,name=iOSBadgeCount,proto3" json:"iOSBadgeCount"`
}

func (x *OfflinePushInfo) Reset() {
	*x = OfflinePushInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OfflinePushInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OfflinePushInfo) ProtoMessage() {}

func (x *OfflinePushInfo) ProtoReflect() protoreflect.Message {
	mi := &file_msg_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OfflinePushInfo.ProtoReflect.Descriptor instead.
func (*OfflinePushInfo) Descriptor() ([]byte, []int) {
	return file_msg_msg_proto_rawDescGZIP(), []int{2}
}

func (x *OfflinePushInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *OfflinePushInfo) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *OfflinePushInfo) GetEx() string {
	if x != nil {
		return x.Ex
	}
	return ""
}

func (x *OfflinePushInfo) GetIOSPushSound() string {
	if x != nil {
		return x.IOSPushSound
	}
	return ""
}

func (x *OfflinePushInfo) GetIOSBadgeCount() bool {
	if x != nil {
		return x.IOSBadgeCount
	}
	return false
}

var File_msg_msg_proto protoreflect.FileDescriptor

var file_msg_msg_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x73, 0x67, 0x2f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x6d, 0x73, 0x67, 0x22, 0x6d, 0x0a, 0x0b, 0x4d, 0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x54,
	0x6f, 0x4d, 0x51, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x26, 0x0a, 0x07, 0x6d,
	0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d,
	0x73, 0x67, 0x2e, 0x4d, 0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x44,
	0x61, 0x74, 0x61, 0x22, 0x98, 0x07, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x76, 0x49,
	0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x63, 0x76, 0x49, 0x44, 0x12,
	0x18, 0x0a, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x12, 0x2a, 0x0a,
	0x10, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x49,
	0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x50,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x49, 0x44, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x46, 0x61, 0x63, 0x65, 0x55,
	0x52, 0x4c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x46, 0x61, 0x63, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x73, 0x67,
	0x46, 0x72, 0x6f, 0x6d, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x46,
	0x72, 0x6f, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x73, 0x65,
	0x71, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x11, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x12, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x4d, 0x73, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3e, 0x0a, 0x0f, 0x6f, 0x66,
	0x66, 0x6c, 0x69, 0x6e, 0x65, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x13, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e,
	0x65, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0f, 0x6f, 0x66, 0x66, 0x6c, 0x69,
	0x6e, 0x65, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x14, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0c, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x6d, 0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x15, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6d, 0x73, 0x67, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x65, 0x78, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x52, 0x65, 0x61, 0x63, 0x74, 0x18,
	0x28, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x52, 0x65, 0x61, 0x63, 0x74, 0x12, 0x32,
	0x0a, 0x14, 0x69, 0x73, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x29, 0x20, 0x01, 0x28, 0x08, 0x52, 0x14, 0x69, 0x73,
	0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x2e, 0x0a, 0x12, 0x6d, 0x73, 0x67, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4d, 0x6f,
	0x64, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x2a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x12,
	0x6d, 0x73, 0x67, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x54, 0x69,
	0x6d, 0x65, 0x1a, 0x3a, 0x0a, 0x0c, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x95,
	0x01, 0x0a, 0x0f, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x0e, 0x0a, 0x02,
	0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x65, 0x78, 0x12, 0x22, 0x0a, 0x0c,
	0x69, 0x4f, 0x53, 0x50, 0x75, 0x73, 0x68, 0x53, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x69, 0x4f, 0x53, 0x50, 0x75, 0x73, 0x68, 0x53, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x24, 0x0a, 0x0d, 0x69, 0x4f, 0x53, 0x42, 0x61, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x4f, 0x53, 0x42, 0x61, 0x64, 0x67,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x4d, 0x53, 0x44, 0x4b, 0x2f, 0x4f,
	0x70, 0x65, 0x6e, 0x2d, 0x49, 0x4d, 0x2d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x6f,
	0x6f, 0x6c, 0x73, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x2f, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_msg_proto_rawDescOnce sync.Once
	file_msg_msg_proto_rawDescData = file_msg_msg_proto_rawDesc
)

func file_msg_msg_proto_rawDescGZIP() []byte {
	file_msg_msg_proto_rawDescOnce.Do(func() {
		file_msg_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_msg_proto_rawDescData)
	})
	return file_msg_msg_proto_rawDescData
}

var file_msg_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_msg_msg_proto_goTypes = []interface{}{
	(*MsgDataToMQ)(nil),     // 0: msg.MsgDataToMQ
	(*MsgData)(nil),         // 1: msg.MsgData
	(*OfflinePushInfo)(nil), // 2: msg.OfflinePushInfo
	nil,                     // 3: msg.MsgData.OptionsEntry
}
var file_msg_msg_proto_depIdxs = []int32{
	1, // 0: msg.MsgDataToMQ.msgData:type_name -> msg.MsgData
	3, // 1: msg.MsgData.options:type_name -> msg.MsgData.OptionsEntry
	2, // 2: msg.MsgData.offlinePushInfo:type_name -> msg.OfflinePushInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_msg_msg_proto_init() }
func file_msg_msg_proto_init() {
	if File_msg_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgDataToMQ); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OfflinePushInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_msg_proto_goTypes,
		DependencyIndexes: file_msg_msg_proto_depIdxs,
		MessageInfos:      file_msg_msg_proto_msgTypes,
	}.Build()
	File_msg_msg_proto = out.File
	file_msg_msg_proto_rawDesc = nil
	file_msg_msg_proto_goTypes = nil
	file_msg_msg_proto_depIdxs = nil
}
