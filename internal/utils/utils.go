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

package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"reflect"
)

func JsonDataList(resp interface{}) []map[string]interface{} {
	var list []proto.Message
	if reflect.TypeOf(resp).Kind() == reflect.Slice {
		s := reflect.ValueOf(resp)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(proto.Message))
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range list {
		m := ProtoToMap(v, false)
		result = append(result, m)
	}
	return result
}

func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}
