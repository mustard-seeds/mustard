package time_util

import (
	"time"
)

func GetCurrentTimeStamp() int64 {
	return time.Now().Unix()
}
func Sleep(s int) {
	time.Sleep(time.Second * time.Duration(s))
}