package storage
import (
    . "mustard/storage/mongo"
    . "mustard/crawl/proto"
)

/*
    Storage API Layer
*/
type StorageEngine interface {
    Save(*CrawlDoc) error
    SaveBatch([]*CrawlDoc) int
    QueryById(docid string, *[]*CrawlDoc) StorageEngine
    QueryByTag(primaryTag string, secondTag []string)
}

func NewStorageEngine() StorageEngine {
    // Decide which engine to use....
    return NewMongoStorageEngine()
}
