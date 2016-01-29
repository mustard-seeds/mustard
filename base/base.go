package base
import "errors"

func CHECK(good bool) {
	if !good {
		panic(errors.New("CHECK fail!"))
	}
}