package babysitter
import (
	"net/http"
	LOG "mustard/base/log"
	"mustard/internal/github.com/gorilla/mux"
	"fmt"
	"mustard/base/string_util"
	"encoding/json"
)

const (
	StatusUiPath = "/statusi"
	StatusUiAPIPath = "/statusi/api/machine"
)

type MonitorInterface interface {
	MonitorReport(*MonitorResult)
}
type MonitorHandleFunc interface {
	MHandleFunc(string, func(w http.ResponseWriter, r *http.Request))
}
type MonitorResult struct {
	info    string
	kv      map[string]string
	machine     string
}
func (mr *MonitorResult)AddString(s string) {
	mr.info = s
}
func (mr *MonitorResult)AddKv(k,v string) {
	mr.kv[k] = v
}

type MonitorServer struct {
	result *MonitorResult
	monitor MonitorInterface
	handleFunc map[string]http.HandlerFunc
}

func (m *MonitorServer)Init() {
	m.result = &MonitorResult{kv:make(map[string]string),
		machine:machineInfoHtml(),
	}
	m.handleFunc = make(map[string]http.HandlerFunc)
}
func (m *MonitorServer)AddMonitor(mi MonitorInterface) {
	m.monitor = mi
}
// custom handlerfunc.
func (m *MonitorServer)AddHandleFunc(path string, f http.HandlerFunc){
	m.handleFunc[path] = f
}
func (m *MonitorServer)StatusiApi(w http.ResponseWriter, r *http.Request) {
	machine := make(map[string]string)
	for k,v := range machineInfo() {
		machine[k] = v
	}
	for k,v := range statusInfo() {
		machine[k] = v
	}
	info,_ := json.Marshal(machine)
	w.Header().Set("Content-Type","application/json")
	w.Write(info)
}
func (m *MonitorServer)Statusi(w http.ResponseWriter, r *http.Request) {
	m.monitor.MonitorReport(m.result)
	w.Header().Set("Server","Golang Server")
	w.WriteHeader(200)
	var infos string
	infos += "<h1>Machine Info</h1>"
	string_util.StringAppendF(&infos,"%s",m.result.machine)
	infos += "<h1>Status Info</h1>"
	string_util.StringAppendF(&infos, "%s", statusInfoHtml())
	infos += "<h1>Application Info</h1>"
	for k,v := range m.result.kv {
		string_util.StringAppendF(&infos, "<key>%s : <value>%s<br>", k, v)
	}
	string_util.StringAppendF(&infos, "<br>%s", m.result.info)
	w.Write([]byte(infos))
}
func (m *MonitorServer)Serve(port int) {
	r := mux.NewRouter()
	r.HandleFunc(StatusUiPath, m.Statusi)
	r.HandleFunc(StatusUiAPIPath, m.StatusiApi)
	LOG.Infof("MonitorServer Serve path %s", StatusUiPath)
	LOG.Infof("MonitorServer Serve path %s", StatusUiAPIPath)
	for k,v := range m.handleFunc {
		r.HandleFunc(k,v)
		LOG.Infof("MonitorServer Serve path %s", k)
	}
	serverAddr := fmt.Sprintf(":%d", port)
	LOG.Infof("Starting Http Monitor at %d", port)
	err := http.ListenAndServe(serverAddr, r)
	if (err != nil) {
		LOG.Fatalf("Http Server Start Fail,%d",port)
	}
}