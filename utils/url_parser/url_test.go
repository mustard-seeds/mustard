package url_parser

import (
	"testing"
)

func TestGetHost(t *testing.T) {
	h := GetHost("http://a.com/index.html")
	if h != "a.com" {
		t.Error("GetHostFail")
	}
}
func TestNormalize1(t *testing.T) {
	u := NormalizeUrl("http://a.com/index.html#xx")
	if u != "http://a.com/index.html" {
		t.Error("Normalize Fail:" + u)
	}
}
