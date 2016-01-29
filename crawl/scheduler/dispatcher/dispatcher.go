package handler

import (
	LOG "mustard/base/log"
	"mustard/base"
	"mustard/base/file"
	"strings"
	"mustard/base/conf"
)
var CONF = conf.Conf
// send crawldoc to target server
// dispatch as:  host/domain/url/random
// remotes:  crawldoc.receivers or configfile.
type CrawlerFeeder struct {
	host string
	port int
}
func (cf *CrawlerFeeder)IsHealthy() bool{
}
func (cf *CrawlerFeeder)IsConnected bool {
}


type CrawlerFeederGroup struct {
	liveFeeders map[uint32]bool
	deadFeeders map[uint32]bool
	feeders map[uint32]*CrawlerFeeder
}

type Dispatcher struct {
}


func (d *Dispatcher)LoadCrawlersFromFile(name string) {
	base.CHECK(file.Exist(name))
	content,_ := file.ReadFileToString(name)
	lines := strings.Split(content, "\n")
	for _,l := range lines {
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		LOG.Info(l)
	}
}
func (d *Dispatcher)init() {

}
