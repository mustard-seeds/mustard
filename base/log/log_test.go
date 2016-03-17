package log

import (
    "testing"
)

func TestLOG(t *testing.T) {
    if VLog(3).level != 3 {
        t.Error("LOG Level set not work.")
    }
    Info("")
    Infof("%d",1)
    Warning("")
    Warningf("")
    Error("")
    Errorf("")
}
