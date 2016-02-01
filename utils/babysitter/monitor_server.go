package babysitter
import (
	"net/http"
	LOG "mustard/base/log"
	"mustard/internal/github.com/gorilla/mux"
	"fmt"
	"mustard/base/string_util"
)


type MonitorInterface interface {
	MonitorReport(*MonitorResult)
}

type MonitorResult struct {
	info    string
	kv      map[string]string
}
func (mr *MonitorResult)AddString(s string) {
	mr.info = s
}
func (mr *MonitorResult)AddKv(k,v string) {
	mr.kv[k] = v
}
func (mr *MonitorResult)machineInfo() string {
	// TODO: mem, cpu, cmd,.....
	return ""
}
type MonitorServer struct {
	result *MonitorResult
	monitor MonitorInterface
}

func (m *MonitorServer)Init() {
	m.result = &MonitorResult{kv:make(map[string]string)}
}
func (m *MonitorServer)AddMonitor(mi MonitorInterface) {
	m.monitor = mi
}
func (m *MonitorServer)Babysitter(w http.ResponseWriter, r *http.Request) {
	m.monitor.MonitorReport(m.result)
	w.Header().Set("Server","Golang Server")
	w.WriteHeader(200)
	var infos string
	string_util.StringAppendF(&infos, "%s", m.result.machineInfo())
	for k,v := range m.result.kv {
		string_util.StringAppendF(&infos, "<key>%s : <value>%s<br>", k, v)
	}
	string_util.StringAppendF(&infos, "<br>%s", m.result.info)
	w.Write([]byte(infos))
}
func (m *MonitorServer)Serve(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/statusi", m.Babysitter)
	LOG.Info("Starting Http Monitor at %d", port)
	serverAddr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(serverAddr, r)
	if (err != nil) {
		LOG.Fatalf("Http Server Start Fail,%d",port)
	}
}