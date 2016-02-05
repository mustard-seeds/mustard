package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "mustard/base/time_util"
    "mustard/base/string_util"
    "fmt"
    "mustard/crawl/fetcher"
)

type FetchHandler struct {
    CrawlHandler
    hostloader *fetcher.HostLoader
    conns   *fetcher.ConnectionPool
    // statistic
    sitequeuefull_num   int
}
func (h *FetchHandler)Status(s *string) {
    h.CrawlHandler.Status(s)
    if h.hostloader != nil && h.conns != nil {
        string_util.StringAppendF(s,  "[(%d/%d)(%d/%d)-%d]",
            h.hostloader.Uim(),
            h.hostloader.Him(),
            h.conns.FreeConnectionNum(),
            h.conns.BusyConnectionNum(),
            h.conns.RecordNum())
    }
}
func (h *FetchHandler)Init() bool {
    LOG.VLog(3).Debugf("FetcherHandler Init")
    h.hostloader = fetcher.NewHostLoader()
    h.conns = fetcher.NewConnectionPool(h.output_chan)
    return true
}
func (h *FetchHandler)Run(p CrawlProcessor) {
    for {
        // non blocking channel
        select {
        case doc := <- h.input_chan:
            LOG.VLog(4).Debugf("Fetcher Handler Get One %s", doc.RequestUrl)
            h.process_num++
            h.accept_num++
            // send to hostload
            e := h.hostloader.Push(doc)
            if e != nil {
                // if push fail, set SITEQUEUEFULL and Output
                LOG.VLog(3).Debugf("SiteQueueFull for %s", doc.Url)
                h.processSiteQueueFul(doc)
                h.sitequeuefull_num++
                h.Output(doc)
            }
        default:
            time_util.Sleep(1)
        }
        fmt.Println()
        h.hostloader.Travel(h.conns.GetCrawlHostMap(), func(doc *proto.CrawlDoc) bool{
            // travel will return the reach time doc. then use connections to fetch.
            // if can not fetch, will return false, then hostloader will not delete it
            return h.conns.Fetch(doc)
        })
        time_util.Sleep(1)
    }
}

func (h *FetchHandler)processSiteQueueFul(doc *proto.CrawlDoc) {
    doc.Code = proto.ReturnType_SITEQUEUEFULLFETCHER
}


// use for create instance from a string
func init() {
    registerCrawlTaskType(&FetchHandler{})
}