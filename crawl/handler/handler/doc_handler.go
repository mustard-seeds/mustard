package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "time"
)

type DocHandler struct {
}

func (doc *DocHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
func (doc *DocHandler)Process(crawlDoc *proto.CrawlDoc) {
    time.Sleep(time.Second * 2)
}
func (doc *DocHandler)Status() {
    LOG.VLog(3).Debug("In DocHandler Status")
}

// use for create instance from a string
func init() {
    registerHandlerType(&DocHandler{})
}