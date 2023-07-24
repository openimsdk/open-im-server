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

package wrapperspb

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
)

func Double(value float64) *DoubleValue {
	return &DoubleValue{Value: value}
}

func Float(value float32) *FloatValue {
	return &FloatValue{Value: value}
}

func Int64(value int64) *Int64Value {
	return &Int64Value{Value: value}
}

func UInt64(value uint64) *UInt64Value {
	return &UInt64Value{Value: value}
}

func Int32(value int32) *Int32Value {
	return &Int32Value{Value: value}
}

func UInt32(value uint32) *UInt32Value {
	return &UInt32Value{Value: value}
}

func Bool(value bool) *BoolValue {
	return &BoolValue{Value: value}
}

func String(value string) *StringValue {
	return &StringValue{Value: value}
}

func Bytes(value []byte) *BytesValue {
	return &BytesValue{Value: value}
}

func DoublePtr(value *float64) *DoubleValue {
	if value == nil {
		return nil
	}
	return &DoubleValue{Value: *value}
}

func FloatPtr(value *float32) *FloatValue {
	if value == nil {
		return nil
	}
	return &FloatValue{Value: *value}
}

func Int64Ptr(value *int64) *Int64Value {
	if value == nil {
		return nil
	}
	return &Int64Value{Value: *value}
}

func UInt64Ptr(value *uint64) *UInt64Value {
	if value == nil {
		return nil
	}
	return &UInt64Value{Value: *value}
}

func Int32Ptr(value *int32) *Int32Value {
	if value == nil {
		return nil
	}
	return &Int32Value{Value: *value}
}

func UInt32Ptr(value *uint32) *UInt32Value {
	if value == nil {
		return nil
	}
	return &UInt32Value{Value: *value}
}

func BoolPtr(value *bool) *BoolValue {
	if value == nil {
		return nil
	}
	return &BoolValue{Value: *value}
}

func StringPtr(value *string) *StringValue {
	if value == nil {
		return nil
	}
	return &StringValue{Value: *value}
}

func BytesPtr(value *[]byte) *BytesValue {
	if value == nil {
		return nil
	}
	return &BytesValue{Value: *value}
}

func (m *DoubleValue) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		return err
	}
	m.Value = value
	return nil
}

func (m *DoubleValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(m.Value, 'f', -1, 64)), nil
}

func (m *FloatValue) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		return err
	}
	m.Value = float32(value)
	return nil
}

func (m *FloatValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(m.Value), 'f', -1, 32)), nil
}

func (m *Int64Value) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return err
	}
	m.Value = value
	return nil
}

func (m *Int64Value) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(m.Value, 10)), nil
}

func (m *UInt64Value) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseUint(string(p), 10, 64)
	if err != nil {
		return err
	}
	m.Value = value
	return nil
}

func (m *UInt64Value) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(m.Value, 10)), nil
}

func (m *Int32Value) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseInt(string(p), 10, 32)
	if err != nil {
		return err
	}
	m.Value = int32(value)
	return nil
}

func (m *Int32Value) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(m.Value), 10)), nil
}

func (m *UInt32Value) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseUint(string(p), 10, 32)
	if err != nil {
		return err
	}
	m.Value = uint32(value)
	return nil
}

func (m *UInt32Value) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(m.Value), 10)), nil
}

func (m *BoolValue) UnmarshalJSON(p []byte) error {
	value, err := strconv.ParseBool(string(p))
	if err != nil {
		return err
	}
	m.Value = value
	return nil
}

func (m *BoolValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatBool(m.Value)), nil
}

func (m *StringValue) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, &m.Value)
}

func (m *StringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (m *BytesValue) UnmarshalJSON(p []byte) error {
	if len(p) < 2 || p[0] != '"' || p[len(p)-1] != '"' {
		return errors.New("invalid bytes value")
	}
	value, err := base64.StdEncoding.DecodeString(string(p[1 : len(p)-1]))
	if err != nil {
		return err
	}
	m.Value = value
	return nil
}

func (m *BytesValue) MarshalJSON() ([]byte, error) {
	return []byte(`"` + base64.StdEncoding.EncodeToString(m.Value) + `"`), nil
}

func (m *DoubleValue) GetValuePtr() *float64 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *FloatValue) GetValuePtr() *float32 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *Int64Value) GetValuePtr() *int64 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *UInt64Value) GetValuePtr() *uint64 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *Int32Value) GetValuePtr() *int32 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *UInt32Value) GetValuePtr() *uint32 {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *BoolValue) GetValuePtr() *bool {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *StringValue) GetValuePtr() *string {
	if m == nil {
		return nil
	}
	return &m.Value
}

func (m *BytesValue) GetValuePtr() *[]byte {
	if m == nil {
		return nil
	}
	return &m.Value
}
