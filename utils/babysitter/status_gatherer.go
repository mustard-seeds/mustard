package babysitter
import (
	"os"
	"mustard/base/string_util"
	"fmt"
)
func machineInfo() map[string]string {
	//TODO pid,cmd, cpunum,total mem, ip port hostname, uptime,
	machine := make(map[string]string)
	var cmd string
	for _,s := range os.Args {
		cmd += s
	}
	machine["cmd"] = cmd
	machine["pid"] = fmt.Sprintf("%d",os.Getpid())
	machine["uid"] = fmt.Sprintf("%d",os.Getuid())
	machine["hostname"],_ = os.Hostname()

	return machine
}
func statusInfo() map[string]string {
	//TODO process mem,cpu,fd, load
	status := make(map[string]string)
	return status
}
func machineInfoHtml() string {
	// just collect one time, information do not change
	machine := machineInfo()
	var info string
	for k,v := range machine {
		string_util.StringAppendF(&info,"<key>%s</key>:<value>%s</value><br>",k,v)
	}
	return info
}
// dynamic, collect when each request
func statusInfoHtml() string {
	return ""
}