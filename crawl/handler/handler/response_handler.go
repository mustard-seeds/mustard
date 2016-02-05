package handler

import (
    "mustard/crawl/proto"
)

type ResponseHandler struct {
    CrawlHandler
    // client cache, could reconnect.
    clients map[string]proto.CrawlServiceClient
}
func (doc *ResponseHandler)Init() bool {
    //TODO response handler, add crawldoc sender
    return false
}
func (doc *ResponseHandler)Accept(crawlDoc *proto.CrawlDoc)bool {
    return true
}
func (doc *ResponseHandler)Process(crawlDoc *proto.CrawlDoc) {
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&ResponseHandler{})
}