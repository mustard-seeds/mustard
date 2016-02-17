package babysitter
import (
	"net/http"
	LOG "mustard/base/log"
	"mustard/internal/github.com/gorilla/mux"
	"fmt"
	"mustard/base/string_util"
	"encoding/json"
	"strings"
	"mustard/base"
	_ "net/http/pprof"
	"net"
)

const (
	StatusUiPath = "/statusi"
	StatusUiAPi =   "/statusi/api"
	StatusUiAPIPath = "/statusi/api/machine"
	StatusUiAPIHealthyPath = "/statusi/api/healthy"
	PprofDebugPath = "/debug/pprof/"
)

type MonitorInterface interface {
	MonitorReport(*MonitorResult)
	MonitorReportHealthy() error
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
// TODO: custom handlerfunc test.... for dispatcher and fetcher add custom api func.
func (m *MonitorServer)AddHandleFunc(path string, f http.HandlerFunc){
	base.CHECK(strings.HasPrefix(path, "/"),"HandleFunc path should start with /")
	path = StatusUiAPi + path
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
func (m *MonitorServer)StatusiHealthyApi(w http.ResponseWriter, r *http.Request) {
	type StatusHealthy struct {
		Healthy bool `json:"Healthy"`
		Reason string `json:"Reason"`
	}
	// json.Marshal only encode Uppercase field in struct...
	info,_ := json.Marshal(StatusHealthy{Healthy:true})
	err := m.monitor.MonitorReportHealthy()
	if err != nil {
		info,_ = json.Marshal(StatusHealthy{Healthy:false,Reason:err.Error()})
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(info)
}
func (m *MonitorServer)Statusi(w http.ResponseWriter, r *http.Request) {
	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	LOG.VLog(3).Debugf("Get Request from %s",ip)
	m.monitor.MonitorReport(m.result)
	w.Header().Set("Server","Golang Statusi Server")
	w.WriteHeader(200)
	var infos string
	infos += "<h1>Machine Info</h1>"
	string_util.StringAppendF(&infos,"%s",m.result.machine)
	infos += "<h1>Status Info</h1>"
	string_util.StringAppendF(&infos, "%s", statusInfoHtml())
	infos += "<h1> Healthy </h1>"
	healthy := fmt.Sprintf("%t", true)
	err := m.monitor.MonitorReportHealthy()
	if err != nil {
		healthy = fmt.Sprintf("%t(%s)", false, err.Error())
	}
	string_util.StringAppendF(&infos, "%s", healthy)
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
	r.HandleFunc(StatusUiAPIHealthyPath, m.StatusiHealthyApi)
	LOG.Infof("MonitorServer Serve path %s", StatusUiPath)
	LOG.Infof("MonitorServer Serve path %s", StatusUiAPIPath)
	LOG.Infof("MonitorServer Serve path %s", StatusUiAPIHealthyPath)
	for k,v := range m.handleFunc {
		r.HandleFunc(k,v)
		LOG.Infof("MonitorServer Serve path %s", k)
	}
	serverAddr := fmt.Sprintf(":%d", port)
	LOG.Infof("Starting Http Monitor at %d", port)

	// AttachProfiler http://stackoverflow.com/questions/19591065/profiling-go-web-application-built-with-gorillas-mux-with-net-http-pprof
	r.PathPrefix(PprofDebugPath).Handler(http.DefaultServeMux)

	err := http.ListenAndServe(serverAddr, r)
	if (err != nil) {
		LOG.Fatalf("Http Server Start Fail,%d",port)
	}
}