package file

import (
	"os"
	"io/ioutil"
	"strings"
//	LOG "mustard/base/log"
	"mustard/base"
	"path/filepath"
	"fmt"
	"mustard/base/conf"
)
var CONF = conf.Conf

func ReadFileToString(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	return string(content),err
}
func Exist(name string) bool {
	_,err := os.Stat(name)
	return !os.IsNotExist(err)
}
// only file, dir is not a file
func FileExist(name string) bool {
	s, err := os.Stat(name)
	return !os.IsNotExist(err) && (!(s != nil && s.IsDir()))
}

func FileLineReader(filename string, comment string, f func(line string)) {
	base.CHECK(Exist(filename), "File %s Not Exist.",filename)
	content,_ := ReadFileToString(filename)
	lines := strings.Split(content,"\n")
	for _,l := range lines {
		if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "#") {
			continue
		}
		f(l)
	}
}

func GetConfFile(s string) string {
	if filepath.IsAbs(s) {
		return s
	}
	cpp := *CONF.ConfPathPrefix
	if strings.HasSuffix(cpp, "mdata") || strings.HasSuffix(cpp, "mdata/") {
		return fmt.Sprintf("%s/%s", *CONF.ConfPathPrefix, s)
	}
	return fmt.Sprintf("%s/mdata/%s", *CONF.ConfPathPrefix, s)
}
