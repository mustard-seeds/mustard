package base
import (
	"errors"
	"fmt"
)

func CHECK(good bool, format string, v ...interface{}) {
	if !good {
		panic(errors.New(fmt.Sprintf("CHECK Fail! " + format, v...)))
	}
}