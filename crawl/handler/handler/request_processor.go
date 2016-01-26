package handler

import (
    "time"
    LOG "mustard/base/log"
    "mustard/crawl/proto"
)

type RequestProcessor struct {
    CrawlHandler
}

func (request *RequestProcessor)Run(p CrawlProcessor) {
    for {
        doc := proto.CrawlDoc{Url:"xxxxurl"}
        time.Sleep(time.Second)
        request.Output(&doc)
        LOG.Info("Send one request")
    }
}
// use for create instance from a string
func init() {
    registerCrawlTaskType(&RequestProcessor{})
}