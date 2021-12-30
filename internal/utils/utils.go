package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func JsonDataList(resp ...interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range resp {
		m := ProtoToMap(v.(proto.Message), false)
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
		EmitDefaults: true,
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
