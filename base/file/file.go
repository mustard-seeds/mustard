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
    name = GetConfFile(name)
    content, err := ioutil.ReadFile(name)
    return string(content),err
}
func WriteStringToFile(content,name string) error {
    return ioutil.WriteFile(name, []byte(content), 0644)
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
    filename = GetConfFile(filename)
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
    /*
        1. if absolute path, return
        2. if ConPathPrefix/s exist. return
        3. else  return ./s
    */
    if filepath.IsAbs(s) {
        return s
    }
    cpp := *CONF.ConfPathPrefix
    realFile := fmt.Sprintf("%s/mdata/%s", *CONF.ConfPathPrefix, s)
    if strings.HasSuffix(cpp, "mdata") || strings.HasSuffix(cpp, "mdata/") {
        realFile = fmt.Sprintf("%s/%s", *CONF.ConfPathPrefix, s)
    }
    if Exist(realFile) {
        return realFile
    } else {
        return s
    }
}
