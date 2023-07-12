// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v4.22.0
// source: wrapperspb/wrapperspb.proto

package wrapperspb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Wrapper message for `double`.
//
// The JSON representation for `DoubleValue` is JSON number.
type DoubleValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The double value.
	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DoubleValue) Reset() {
	*x = DoubleValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoubleValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoubleValue) ProtoMessage() {}

func (x *DoubleValue) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoubleValue.ProtoReflect.Descriptor instead.
func (*DoubleValue) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{0}
}

func (x *DoubleValue) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `float`.
//
// The JSON representation for `FloatValue` is JSON number.
type FloatValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The float value.
	Value float32 `protobuf:"fixed32,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FloatValue) Reset() {
	*x = FloatValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FloatValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FloatValue) ProtoMessage() {}

func (x *FloatValue) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FloatValue.ProtoReflect.Descriptor instead.
func (*FloatValue) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{1}
}

func (x *FloatValue) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `int64`.
//
// The JSON representation for `Int64Value` is JSON string.
type Int64Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The int64 value.
	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Int64Value) Reset() {
	*x = Int64Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int64Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int64Value) ProtoMessage() {}

func (x *Int64Value) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int64Value.ProtoReflect.Descriptor instead.
func (*Int64Value) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{2}
}

func (x *Int64Value) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `uint64`.
//
// The JSON representation for `UInt64Value` is JSON string.
type UInt64Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The uint64 value.
	Value uint64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *UInt64Value) Reset() {
	*x = UInt64Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UInt64Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UInt64Value) ProtoMessage() {}

func (x *UInt64Value) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UInt64Value.ProtoReflect.Descriptor instead.
func (*UInt64Value) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{3}
}

func (x *UInt64Value) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `int32`.
//
// The JSON representation for `Int32Value` is JSON number.
type Int32Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The int32 value.
	Value int32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Int32Value) Reset() {
	*x = Int32Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int32Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int32Value) ProtoMessage() {}

func (x *Int32Value) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int32Value.ProtoReflect.Descriptor instead.
func (*Int32Value) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{4}
}

func (x *Int32Value) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `uint32`.
//
// The JSON representation for `UInt32Value` is JSON number.
type UInt32Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The uint32 value.
	Value uint32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *UInt32Value) Reset() {
	*x = UInt32Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UInt32Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UInt32Value) ProtoMessage() {}

func (x *UInt32Value) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UInt32Value.ProtoReflect.Descriptor instead.
func (*UInt32Value) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{5}
}

func (x *UInt32Value) GetValue() uint32 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Wrapper message for `bool`.
//
// The JSON representation for `BoolValue` is JSON `true` and `false`.
type BoolValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The bool value.
	Value bool `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BoolValue) Reset() {
	*x = BoolValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolValue) ProtoMessage() {}

func (x *BoolValue) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolValue.ProtoReflect.Descriptor instead.
func (*BoolValue) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{6}
}

func (x *BoolValue) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

// Wrapper message for `string`.
//
// The JSON representation for `StringValue` is JSON string.
type StringValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The string value.
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringValue) Reset() {
	*x = StringValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringValue) ProtoMessage() {}

func (x *StringValue) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringValue.ProtoReflect.Descriptor instead.
func (*StringValue) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{7}
}

func (x *StringValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Wrapper message for `bytes`.
//
// The JSON representation for `BytesValue` is JSON string.
type BytesValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The bytes value.
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BytesValue) Reset() {
	*x = BytesValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wrapperspb_wrapperspb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BytesValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BytesValue) ProtoMessage() {}

func (x *BytesValue) ProtoReflect() protoreflect.Message {
	mi := &file_wrapperspb_wrapperspb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BytesValue.ProtoReflect.Descriptor instead.
func (*BytesValue) Descriptor() ([]byte, []int) {
	return file_wrapperspb_wrapperspb_proto_rawDescGZIP(), []int{8}
}

func (x *BytesValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_wrapperspb_wrapperspb_proto protoreflect.FileDescriptor

var file_wrapperspb_wrapperspb_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x70, 0x62, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x4f,
	0x70, 0x65, 0x6e, 0x49, 0x4d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x22, 0x23, 0x0a, 0x0b, 0x44, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x46, 0x6c, 0x6f,
	0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a,
	0x0a, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x23, 0x0a, 0x0b, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x55, 0x49,
	0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x21, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x3a, 0x5a, 0x38, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x4d,
	0x53, 0x44, 0x4b, 0x2f, 0x4f, 0x70, 0x65, 0x6e, 0x2d, 0x49, 0x4d, 0x2d, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wrapperspb_wrapperspb_proto_rawDescOnce sync.Once
	file_wrapperspb_wrapperspb_proto_rawDescData = file_wrapperspb_wrapperspb_proto_rawDesc
)

func file_wrapperspb_wrapperspb_proto_rawDescGZIP() []byte {
	file_wrapperspb_wrapperspb_proto_rawDescOnce.Do(func() {
		file_wrapperspb_wrapperspb_proto_rawDescData = protoimpl.X.CompressGZIP(file_wrapperspb_wrapperspb_proto_rawDescData)
	})
	return file_wrapperspb_wrapperspb_proto_rawDescData
}

var file_wrapperspb_wrapperspb_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_wrapperspb_wrapperspb_proto_goTypes = []interface{}{
	(*DoubleValue)(nil), // 0: OpenIMServer.protobuf.DoubleValue
	(*FloatValue)(nil),  // 1: OpenIMServer.protobuf.FloatValue
	(*Int64Value)(nil),  // 2: OpenIMServer.protobuf.Int64Value
	(*UInt64Value)(nil), // 3: OpenIMServer.protobuf.UInt64Value
	(*Int32Value)(nil),  // 4: OpenIMServer.protobuf.Int32Value
	(*UInt32Value)(nil), // 5: OpenIMServer.protobuf.UInt32Value
	(*BoolValue)(nil),   // 6: OpenIMServer.protobuf.BoolValue
	(*StringValue)(nil), // 7: OpenIMServer.protobuf.StringValue
	(*BytesValue)(nil),  // 8: OpenIMServer.protobuf.BytesValue
}
var file_wrapperspb_wrapperspb_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_wrapperspb_wrapperspb_proto_init() }
func file_wrapperspb_wrapperspb_proto_init() {
	if File_wrapperspb_wrapperspb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wrapperspb_wrapperspb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoubleValue); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FloatValue); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int64Value); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UInt64Value); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int32Value); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UInt32Value); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolValue); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringValue); i {
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
		file_wrapperspb_wrapperspb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BytesValue); i {
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
			RawDescriptor: file_wrapperspb_wrapperspb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wrapperspb_wrapperspb_proto_goTypes,
		DependencyIndexes: file_wrapperspb_wrapperspb_proto_depIdxs,
		MessageInfos:      file_wrapperspb_wrapperspb_proto_msgTypes,
	}.Build()
	File_wrapperspb_wrapperspb_proto = out.File
	file_wrapperspb_wrapperspb_proto_rawDesc = nil
	file_wrapperspb_wrapperspb_proto_goTypes = nil
	file_wrapperspb_wrapperspb_proto_depIdxs = nil
}
