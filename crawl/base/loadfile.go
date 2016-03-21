package base

import (
	"mustard/base/conf"
	"mustard/base/file"
	LOG "mustard/base/log"
	"mustard/base/time_util"
	"strings"
)

var CONF = conf.Conf

func LoadConfigWithTwoField(name, filename, splitS string, last_load_time *int64) (map[string]string, bool) {
	if time_util.GetCurrentTimeStamp()-*last_load_time < int64(*CONF.Crawler.ConfigFileReloadInterval) {
		return nil, false
	}
	filename = file.GetConfFile(filename)
	result := make(map[string]string)
	*last_load_time = time_util.GetCurrentTimeStamp()
	LOG.Infof("Load Config %s", filename)
	file.FileLineReader(filename, "#", func(line string) {
		addr := strings.Split(line, splitS)
		if len(addr) != 2 {
			LOG.Errorf("%s Load Config Format Error, %s : %s", name, filename, line)
			return
		}
		addr0 := strings.TrimSpace(addr[0])
		result[addr0] = addr[1]
		LOG.VLog(4).Debugf("Load %s  %s : %s", name, addr0, addr[1])
	})
	return result, true
}
