package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "mustard/base/time_util"
    "mustard/base/string_util"
)

type FetchHandler struct {
    CrawlHandler
    hostloader *HostLoader
    conns   *ConnectionPool
    // statistic
    sitequeuefull_num   int
}
func (h *FetchHandler)Status(s *string) {
    h.CrawlHandler.Status(s)
    if h.hostloader != nil && h.conns != nil {
        string_util.StringAppendF(s,  "[(%d-%d)-(%d-%d)-%d]",
            h.hostloader.Him(),
            h.hostloader.Uim(),
            h.conns.BusyConnectionNum(),
            h.conns.FreeConnectionNum(),
            len(h.conns.record))
    }
}
func (h *FetchHandler)Init() {
    h.hostloader = &HostLoader{
        hostMap: make(map[string]*HostLoadQueue),
    }
    h.conns = &ConnectionPool{
        record : make(map[string]int64),
        hold : make(map[string]bool),
        busy : make(map[*Connection]bool),
        output_chan : h.output_chan,
    }
}
func (h *FetchHandler)Run(p CrawlProcessor) {
    h.Init()
    for {
        // non blocking channel
        select {
        case doc := <- h.input_chan:
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
        h.hostloader.Travel(h.conns.GetCrawlHostMap(), func(doc *proto.CrawlDoc) bool{
            // travel will return the reach time doc. then use connections to fetch.
            // if can not fetch, will return false, then hostloader will not delete it
            return h.conns.Fetch(doc)
        })
        time_util.Sleep(1)
    }
    h.crawlDoc = <- h.input_chan
}

func (h *FetchHandler)processSiteQueueFul(doc *proto.CrawlDoc) {
    doc.Code = proto.ReturnType_SITEQUEUEFULL
}


// use for create instance from a string
func init() {
    registerCrawlTaskType(&FetchHandler{})
}