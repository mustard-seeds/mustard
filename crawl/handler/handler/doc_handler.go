package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
)

type DocHandler struct {
}

func (doc *DocHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
func (doc *DocHandler)Process(crawlDoc *proto.CrawlDoc) {

}
func (doc *DocHandler)Status() {
    LOG.VLog(3).Debug("In DocHandler Status")
}

// use for create instance from a string
func init() {
    registerType((*DocHandler)(nil))
}