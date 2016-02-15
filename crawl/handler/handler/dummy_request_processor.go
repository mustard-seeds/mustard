package handler

import (
    "time"
    LOG "mustard/base/log"
    "mustard/crawl/proto"
)

type DummyRequestProcessor struct {
    CrawlHandler
}

func (request *DummyRequestProcessor)Run(p CrawlProcessor) {
    for {
        doc := new(proto.CrawlDoc)
        doc.Url = "http://www.a.com/index.html"
        doc.CrawlParam = new(proto.CrawlParam)
        doc.CrawlParam.FetchHint = new(proto.FetchHint)
        doc.CrawlParam.FetchHint.Host = "a.com"
        doc.CrawlParam.Hostload = 5
        doc.CrawlRecord = new(proto.CrawlRecord)
        time.Sleep(time.Second)
        request.Output(doc)
        LOG.Info("Send one request")
    }
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&DummyRequestProcessor{})
}