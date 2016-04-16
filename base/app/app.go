package app

import (
	"runtime"
	LOG "mustard/base/log"
)
const (
	kDefaultMaxProc = 8
)

func init() {
	maxProc := runtime.NumCPU()
	if maxProc > kDefaultMaxProc {
		maxProc = kDefaultMaxProc
	}
	currentProc := runtime.GOMAXPROCS(maxProc)
	LOG.Infof("GOMAXPROC %d ==>> %d",currentProc,maxProc)
}
