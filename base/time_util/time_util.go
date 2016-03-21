package time_util

import (
	"time"
)

func GetCurrentTimeStamp() int64 {
	return time.Now().Unix()
}
func GetTimeInMs() int64 {
	return time.Now().UnixNano() / int64(1000000)
}
func GetReadableTimeNow() string {
	return GetReadableTime(GetCurrentTimeStamp())
}

// input is timestamp
func GetReadableTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04 MST")
}
func Sleep(s int) {
	time.Sleep(time.Second * time.Duration(s))
}
func Usleep(n int) {
	time.Sleep(time.Microsecond * time.Duration(n))
}
