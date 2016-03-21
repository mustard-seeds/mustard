package time_util

import (
	"testing"
)

func TestGetReadableTime(t *testing.T) {
	at := GetReadableTime(1454387433)
	if at != "2016-02-02 12:30 CST" {
		t.Error("GetReadableTime Error")
	}
}
