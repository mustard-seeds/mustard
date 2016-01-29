package main

import (
	"mustard/crawl/handler"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	c := handler.CrawlHandlerController{}
	c.InitCrawlService()
	for {
		c.PrintStatus()
		time.Sleep(time.Second*10)
	}
}

