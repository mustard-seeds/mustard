package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "time"
    "reflect"
)

type DocHandler struct {
    CrawlHandler
}

func (doc *DocHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
func (doc *DocHandler)Process(crawlDoc *proto.CrawlDoc) {
    LOG.VLog(3).Debugf("[%s]Process One Doc %s ",
        reflect.Indirect(reflect.ValueOf(doc)).Type().Name(),
        crawlDoc.Url)
    time.Sleep(time.Second * 2)
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&DocHandler{})
}