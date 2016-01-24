package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
)

type StorageHandler struct {

}
func (doc *StorageHandler)Accept(crawlDoc *proto.CrawlDoc)bool {
    return true
}
func (doc *StorageHandler)Process(crawlDoc *proto.CrawlDoc) {

}
func (doc *StorageHandler)Status() {
    LOG.VLog(3).Debug("In StorageHandler Status")
}

// use for create instance from a string
func init() {
    registerType((*StorageHandler)(nil),)
}