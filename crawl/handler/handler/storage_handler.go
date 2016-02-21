package handler

import (
    "mustard/crawl/proto"
    . "mustard/storage"
    "sync"
    "mustard/base/time_util"
    LOG "mustard/base/log"
    "mustard/base/proto_util"
    "mustard/base/string_util"
)

const (
    kSaveBatchSize = 2
    kSaveBatchInterval = 300  // second
)
type StorageHandler struct {
    CrawlHandler
    docs []*proto.CrawlDoc
    sync.RWMutex
    last_db_time int64
}

func (handler *StorageHandler)saveDocs() {
    handler.Lock()
    defer handler.Unlock()
    if len(handler.docs) == 0 {
        return
    }
    if len(handler.docs) > kSaveBatchSize ||
        time_util.GetCurrentTimeStamp() - handler.last_db_time > kSaveBatchInterval {
        t1 := time_util.GetTimeInMs()
        num,err := STORAGE_ENGINE_IMPL.WithDb(*CONF.Crawler.ContentDbName).
                        WithTable(*CONF.Crawler.ContentDbTable).
                        SaveBatch(handler.docs)
        if err != nil {
            LOG.VLog(2).Debugf("Save Content to db error %s",err.Error())
        }
        handler.docs = nil
        handler.last_db_time = time_util.GetCurrentTimeStamp()
        LOG.VLog(3).Debugf("Flush %d using time %d ms.",num,time_util.GetTimeInMs() - t1)
    }
}
func (handler *StorageHandler)DBThread() {
    for true {
        handler.saveDocs()
        time_util.Sleep(10)
    }
}
func (handler *StorageHandler)Init() bool {
    STORAGE_ENGINE_IMPL.Init(*CONF.Crawler.ContentDBServers)
    go handler.DBThread()
    return true
}
func (handler *StorageHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}
// save doc to content db
func (handler *StorageHandler)Process(crawlDoc *proto.CrawlDoc) {
    handler.Lock()
    defer handler.Unlock()
    // deepcopy crawldoc
    docMsg := proto_util.FromProtoToString(crawlDoc)
    var newDoc proto.CrawlDoc
    proto_util.FromStringToProto(docMsg, &newDoc)
    // compress and save
    compressContent,err := string_util.Compress(newDoc.Content)
    if err == nil {
        newDoc.Content = compressContent
        newDoc.ContentCompressed = true
    } else {
        LOG.VLog(2).Debugf("Compress Error url:%s,docid:%s,error:%s",
                    newDoc.Url,newDoc.Docid,err.Error())
    }
    handler.docs = append(handler.docs, &newDoc)
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&StorageHandler{})
}
