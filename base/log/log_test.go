package log

import (
    "testing"
)

func TestLOG(t *testing.T) {
    VLog(3).Debug("")
    if _log.level != 3 {
        t.Error("LOG Level set not work.")
    }
    Info("")
    Infof("%d",1)
    Warning("")
    Warningf("")
    Error("")
    Errorf("")
}
