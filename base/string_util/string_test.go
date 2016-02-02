package string_util
import (
	"testing"
)
func TestCompressUnCompress(t *testing.T) {
	str := "this is origin text which need to be compressed"
	cstr,err := Compress(str)
	if err != nil {
		t.Error("Compress throw error")
	}
	rstr,err2 := Uncompress(cstr)
	if err2 != nil {
		t.Error("UnCompress throw error")
	}
	if str != rstr {
		t.Error("Compress -> UnCompress not equal")
	}
}
