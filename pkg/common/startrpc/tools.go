package startrpc

import (
	"reflect"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func getConfig[T any](value reflect.Value) *T {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() == reflect.Struct {
		num := value.NumField()
		for i := 0; i < num; i++ {
			field := value.Field(i)
			for field.Kind() == reflect.Pointer {
				field = field.Elem()
			}
			if field.Kind() == reflect.Struct {
				if elem, ok := field.Interface().(T); ok {
					return &elem
				}
				if elem := getConfig[T](field); elem != nil {
					return elem
				}
			}
		}
	}
	return nil
}

func getConfigRpcMaxRequestBody(value reflect.Value) *conf.MaxRequestBody {
	return getConfig[conf.MaxRequestBody](value)
}

func getConfigShare(value reflect.Value) *conf.Share {
	return getConfig[conf.Share](value)
}
