package time_util

import (
    "time"
)

func GetCurrentTimeStamp() int64 {
    return time.Now().Unix()
}
func GetTimeInMs() int64 {
    return time.Now().UnixNano() / int64(1000)
}
func Sleep(s int) {
    time.Sleep(time.Second * time.Duration(s))
}