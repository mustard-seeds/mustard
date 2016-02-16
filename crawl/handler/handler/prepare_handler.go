package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "mustard/crawl/base"
)

type PrepareHandler struct {
    CrawlHandler
}

func (doc *PrepareHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
func (doc *PrepareHandler)Process(crawlDoc *proto.CrawlDoc) {
    LOG.VLog(4).Debugf("\n%s", base.DumpCrawlDoc(crawlDoc))
    LOG.VLog(4).Debugf("Content:\n%U", crawlDoc.Content)
    //TODO format, encode... exchange
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&PrepareHandler{})
}