package conf

import (
    "testing"
)

func TestConf(t *testing.T) {
    var c = Conf
    if *c.LogV < 0 {
        t.Error("LOG set Level < 0")
    }
}
