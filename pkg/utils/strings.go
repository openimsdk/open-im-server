/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/4/8 15:09).
 */
package utils

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

// transfer int to string
func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// transfer string to int
func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}

// transfer string to int64
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

// transfer string to int32
func StringToInt32(i string) int32 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return int32(j)
}

// transfer int32 to string
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// transfer unit32 to string
func Uint32ToString(i uint32) string {
	return strconv.FormatInt(int64(i), 10)
}

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {
	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false
}

// contain int32 or not
func IsContainInt32(target int32, List []int32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}

// contain int or not
func IsContainInt(target int, List []int) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}

// transfer array to string array
func InterfaceArrayToStringArray(data []interface{}) (i []string) {
	for _, param := range data {
		i = append(i, param.(string))
	}
	return i
}

// transfer struct to json string
func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// transfer struct to jsonBytes
func StructToJsonBytes(param interface{}) []byte {
	dataType, _ := json.Marshal(param)
	return dataType
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
}

// get message ID
func GetMsgID(sendID string) string {
	t := int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}

// transfer int64 to string
func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// remove duplicate elment
func RemoveDuplicateElement(idList []string) []string {
	result := make([]string, 0, len(idList))
	temp := map[string]struct{}{}
	for _, item := range idList {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// remove duplicate
func RemoveDuplicate[T comparable](arr []T) []T {
	result := make([]T, 0, len(arr))
	temp := map[T]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// judege is or not a duplicated string clice
func IsDuplicateStringSlice(arr []string) bool {
	t := make(map[string]struct{})
	for _, s := range arr {
		if _, ok := t[s]; ok {
			return true
		}
		t[s] = struct{}{}
	}
	return false
}
