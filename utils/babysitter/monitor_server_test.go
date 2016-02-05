package babysitter

import (
	"testing"
	"net/http"
)

func TestMonitorResult(t *testing.T) {
	mr := &MonitorResult{
		kv:make(map[string]string),
	}
	mr.AddString("X")
	mr.AddString("XX")
	if mr.info != "XX" {
		t.Error("MonitorResult Add String not overwrite")
	}
	mr.AddKv("1","1")
	mr.AddKv("1","1")
	if len(mr.kv) != 1 {
		t.Error("MonitorResult AddKv not update")
	}
}
func TestMonitorServer(t *testing.T) {
	ms := &MonitorServer{
		handleFunc:make(map[string]http.HandlerFunc),
	}
	ms.AddHandleFunc("/a",nil)
	_,p := ms.handleFunc[StatusUiAPi+"/a"]
	if p == false {
		t.Error("Custom HandleFunc not use StatusUiApi prefix")
	}
}
