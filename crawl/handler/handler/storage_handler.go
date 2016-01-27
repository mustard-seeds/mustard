package handler

import (
    "mustard/crawl/proto"
)

type StorageHandler struct {
    CrawlHandler
}
func (doc *StorageHandler)Accept(crawlDoc *proto.CrawlDoc)bool {
    return true
}
func (doc *StorageHandler)Process(crawlDoc *proto.CrawlDoc) {
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&StorageHandler{})
}