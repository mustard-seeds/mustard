package handler

import (
    "mustard/crawl/proto"
)

type StorageHandler struct {
    CrawlHandler
}

func (handler *StorageHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
// save doc to content db
func (handler *StorageHandler)Process(crawlDoc *proto.CrawlDoc) {
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&StorageHandler{})
}
