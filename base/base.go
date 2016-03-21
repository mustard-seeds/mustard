package base

import (
	"errors"
	"fmt"
)

func CHECK(good bool, format string, v ...interface{}) {
	if !good {
		panic(errors.New(fmt.Sprintf("CHECK Fail! "+format, v...)))
	}
}
func CHECKERROR(e error, format string, v ...interface{}) {
	if e != nil {
		panic(errors.New(fmt.Sprintf("CHECK ERROR FAIL, Error(%s)  ", e.Error()) +
			fmt.Sprintf(format, v...)))
	}
}
