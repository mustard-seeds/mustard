package handler

import (
    "time"
    "mustard/crawl/proto"
)

type RequestProcessor struct {
    output_chan  chan<- *proto.CrawlDoc
}

func (request *RequestProcessor)GetHandler() CrawlHandler {
    return (CrawlHandler)(nil)
}
func (request *RequestProcessor)SetHandler(CrawlHandler) {
}
func (request *RequestProcessor)Run() {
    for {
        doc := proto.CrawlDoc{Url:"xxxxurl"}
        time.Sleep(time.Second)
        request.Output(&doc)
    }
}
func (request *RequestProcessor)Output(doc *proto.CrawlDoc) {
    if request.output_chan != nil {
        request.output_chan <- doc
    }
}

func (request *RequestProcessor)GetInputChan() <-chan *proto.CrawlDoc {
    return nil
}
func (request *RequestProcessor)GetOutputChan() chan<- *proto.CrawlDoc {
    return request.output_chan
}
func (request *RequestProcessor)SetInputChan(out chan<-*proto.CrawlDoc) {
}
func (request *RequestProcessor)SetOutputChan(out chan<-*proto.CrawlDoc) {
    request.output_chan = out
}

// use for create instance from a string
func init() {
    registerType((*RequestProcessor)(nil))
}