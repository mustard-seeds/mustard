package storage

import (
    "mustard/crawl/proto"
)

/*
    Storage API Layer
*/
type StorageEngine interface {
    // crawldoc interface
    Save(*proto.CrawlDoc) error
    SaveBatch([]*proto.CrawlDoc) int
    QueryById(docid string, *[]*proto.CrawlDoc) error
    //TODO  what type api here???   filter for which field???
}
