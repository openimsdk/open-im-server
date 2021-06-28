/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/2/22 11:52).
 */
package log

import (
	"strconv"
	"time"
)

const (
	TimeOffset = 8 * 3600  //8个小时的偏移量
	HalfOffset = 12 * 3600 //半天的小时偏移量
)

//获取当前的时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

//获取当天0点的时间戳
func GetCurDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - TimeOffset
}

//获取当天12点的时间戳
func GetCurDayHalfTimestamp() int64 {
	return GetCurDayZeroTimestamp() + HalfOffset

}

//获取当天0点格式化时间，格式为"2006-01-02_00-00-00"
func GetCurDayZeroTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp(), 0).Format("2006-01-02_15-04-05")
}

//获取当天12点格式化时间，格式为"2006-01-02_12-00-00"
func GetCurDayHalfTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp()+HalfOffset, 0).Format("2006-01-02_15-04-05")
}
func GetTimeStampByFormat(datetime string) string {
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64
	return strconv.FormatInt(timestamp, 10)
}

func TimeStringFormatTimeUnix(timeFormat string, timeSrc string) int64 {
	tm, _ := time.Parse(timeFormat, timeSrc)
	return tm.Unix()
}
