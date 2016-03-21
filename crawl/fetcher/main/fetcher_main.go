package main

import (
	"mustard/base/conf"
	"mustard/crawl/handler"
	"mustard/utils/babysitter"
	"runtime"
	"time"
)

var CONF = conf.Conf

func main() {
	runtime.GOMAXPROCS(2)
	c := handler.CrawlHandlerController{}
	c.InitCrawlService()

	var http_server babysitter.MonitorServer
	http_server.Init()

	http_server.AddMonitor(&c)
	go http_server.Serve(*CONF.Crawler.HttpPort)

	for {
		c.PrintStatus()
		time.Sleep(time.Second * 10)
	}
}
