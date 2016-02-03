package file

import (
	"testing"
)
func TestGetConfFile(t *testing.T) {
	fname := GetConfFile("/usr/local/bin/go")
	if fname != "/usr/local/bin/go" {
		t.Error("GetConfFile Fail. absolute path.")
	}
}
