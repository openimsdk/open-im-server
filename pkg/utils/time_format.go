/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/2/22 11:52).
 */
package utils

import (
	"strconv"
	"time"
)

const (
	TimeOffset = 8 * 3600  //8 hour offset
	HalfOffset = 12 * 3600 //Half-day hourly offset
)

//Get the current timestamp by Second
func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

//Convert timestamp to time.Time type
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

//Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}
func UnixMillSecondToTime(millSecond int64) time.Time {
	return time.Unix(0, millSecond*1e6)
}

//Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

//Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

//Get the timestamp at 0 o'clock of the day
func GetCurDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - TimeOffset
}

//Get the timestamp at 12 o'clock on the day
func GetCurDayHalfTimestamp() int64 {
	return GetCurDayZeroTimestamp() + HalfOffset

}

//Get the formatted time at 0 o'clock of the day, the format is "2006-01-02_00-00-00"
func GetCurDayZeroTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp(), 0).Format("2006-01-02_15-04-05")
}

//Get the formatted time at 12 o'clock of the day, the format is "2006-01-02_12-00-00"
func GetCurDayHalfTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp()+HalfOffset, 0).Format("2006-01-02_15-04-05")
}
func GetTimeStampByFormat(datetime string) string {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return strconv.FormatInt(timestamp, 10)
}

func TimeStringFormatTimeUnix(timeFormat string, timeSrc string) int64 {
	tm, _ := time.Parse(timeFormat, timeSrc)
	return tm.Unix()
}

func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}
