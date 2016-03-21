package string_util

import (
	"fmt"
	"strings"
)

func Purify(s string, dirty ...string) string {
	n := s
	for _, d := range dirty {
		n = strings.Replace(n, d, "", -1)
	}
	return n
}

func IsEmpty(s string) bool {
	return s == ""
}

func StringAppendF(s *string, format string, a ...interface{}) {
	app := fmt.Sprintf(format, a...)
	*s = *s + app
}
