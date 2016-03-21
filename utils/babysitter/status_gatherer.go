package babysitter

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"mustard/base/time_util"
	"mustard/base/string_util"
	"mustard/internal/github.com/c9s/goprocinfo/linux"
)

func machineInfo() map[string]string {
	//TODO pid,cmd, cpunum,total mem, ip port hostname, uptime,
	machine := make(map[string]string)
	machine["cmd"] = strings.Join(os.Args, " ")
	machine["pid"] = fmt.Sprintf("%d", os.Getpid())
	machine["uid"] = fmt.Sprintf("%d", os.Getuid())
	machine["hostname"], _ = os.Hostname()
	machine["StartAt"] = time_util.GetReadableTimeNow()
	return machine
}
func statusInfo() map[string]string {
	// TODO process mem,cpu,fd, load
	// https://github.com/c9s/goprocinfo
	pid := os.Getpid()
	processPath:= filepath.Join("/proc", strconv.Itoa(pid))
	status := make(map[string]string)
	mem,err1 := linux.ReadMemInfo("/proc/meminfo")
	if err1 == nil {
		status["MemTotal"] = strconv.Itoa(int(mem.MemTotal/1000))
		status["MemFree"] = strconv.Itoa(int(mem.MemFree + mem.Cached + mem.Buffers)/1000)
	}
	processMem, err2 := linux.ReadProcessStatus(filepath.Join(processPath, "status"))
	if err2 == nil {
		status["VmRSS"] = strconv.Itoa(int(processMem.VmRSS/1000))
		status["VmSize"] = strconv.Itoa(int(processMem.VmSize/1000))
	}
	if err1 == nil && err2 == nil {
		status["MemUse"] = strconv.FormatFloat(float64(processMem.VmRSS)/float64(mem.MemTotal), 'f', -1, 64)
	}
	files, err3 := ioutil.ReadDir(filepath.Join(processPath, "fd"))
	if err3 == nil {
		status["FDS"] = strconv.Itoa(len(files))
	}
	// TODO CPU Usage http://stackoverflow.com/questions/16726779/how-do-i-get-the-total-cpu-usage-of-an-application-from-proc-pid-stat

	return status
}
func machineInfoHtml() string {
	// just collect one time, information do not change
	machine := machineInfo()
	var info string
	for k, v := range machine {
		string_util.StringAppendF(&info, "<key>%s</key>:<value>%s</value><br>", k, v)
	}
	return info
}

// dynamic, collect when each request
func statusInfoHtml() string {
	return ""
}
