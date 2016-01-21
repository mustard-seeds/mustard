package string_util
import (
	"strings"
)

func Purify(s string, dirty ...string) string {
	n := s
	for _,d := range dirty {
		n =strings.Replace(n, d, "", -1)
	}
	return n
}