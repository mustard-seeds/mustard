package file

import (
	"os"
	"io/ioutil"
//	LOG "mustard/base/log"
)

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

