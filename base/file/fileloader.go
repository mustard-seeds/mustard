package file

import (
	"strings"
	"strconv"
	"fmt"
	. "mustard/base"
	"mustard/base/time_util"
	LOG "mustard/base/log"
)

type ConfigLoader struct {
	Last_load_time int64
	Last_load_version int

	Reload_interval int64
	Last_version int

	Version_file string
}
func (c *ConfigLoader)ShouldReloadConfig() bool {
	// only call once...
	if time_util.GetCurrentTimeStamp() - c.Last_load_time > c.Reload_interval {
		c.Last_load_time = time_util.GetCurrentTimeStamp()
		newVersion := c.GetVersionFromFile()
		if c.Last_version < newVersion {
			c.Last_version = newVersion
			return true
		}
	}
	return false
}
func (c *ConfigLoader)GetVersionFromFile() int {
	currentVersion := -1
	FileLineReaderSoftly(c.Version_file, "#", func(line string){
		value, e := strconv.Atoi(line)
		if e != nil {
			LOG.Errorf("GerVersionFail %s", line)
		} else {
			if currentVersion == -1 {
				currentVersion = value
			}
		}
	})
	if currentVersion < 0 {
		LOG.Errorf("GerVersion Fail %s", c.Version_file)
	}
	return currentVersion
}
func (c *ConfigLoader)SetConfigVersionFile(filename string) *ConfigLoader {
	c.Version_file = GetConfFile(filename)
	CHECK(Exist(c.Version_file), "ConfigLoader %s not exist", c.Version_file)
	c.Last_load_time = -1
	return c
}
func (c *ConfigLoader)SetReloadInterval(interval int64) *ConfigLoader {
	c.Reload_interval = interval
	CHECK(c.Reload_interval > 0, "interval should bigger then 0")
	return c
}
func (c *ConfigLoader)NextLoadInfo() string {
	remain := c.Reload_interval - time_util.GetCurrentTimeStamp() + c.Last_load_time
	if remain < 0 {
		remain = 0
	}
	return fmt.Sprintf("Remain:%d,Version:%d", remain, c.Last_version)
}
func (c *ConfigLoader)LoadConfigWithTwoField(name, filename, splitS string) map[string]string {
	c.Last_load_time = time_util.GetCurrentTimeStamp()
	filename = GetConfFile(filename)
	result := make(map[string]string)
	LOG.Infof("Load Config %s", filename)
	FileLineReader(filename, "#", func(line string) {
		addr := strings.Split(line, splitS)
		if len(addr) != 2 {
			LOG.Errorf("%s Load Config Format Error, %s : %s", name, filename, line)
			return
		}
		addr0 := strings.TrimSpace(addr[0])
		result[addr0] = strings.TrimSpace(addr[1])
		LOG.VLog(5).Debugf("Load %s  %s : %s", name, addr0, addr[1])
	})
	return result
}