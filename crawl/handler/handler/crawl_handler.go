package handler

import (
    "mustard/base/time_util"
    "reflect"
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
)

var CONF = conf.Conf

type CrawlTask interface {
    PeriodicTask(p CrawlProcessor)
    Run(p CrawlProcessor)
    Output(doc *proto.CrawlDoc)
    SetInputChan(in <-chan*proto.CrawlDoc)
    SetOutputChan(out chan<-*proto.CrawlDoc)
    GetInputChan() <-chan *proto.CrawlDoc
    GetOutputChan() chan<- *proto.CrawlDoc
}

type CrawlProcessor interface {
    Status()
    Process(*proto.CrawlDoc)
    Accept(*proto.CrawlDoc) bool
}

type CrawlHandler struct {
    input_chan <-chan *proto.CrawlDoc
    output_chan  chan<- *proto.CrawlDoc
    crawlDoc *proto.CrawlDoc

    // statistic
    process_num         int64
    accept_num          int64
    max_process_time    int64
    avg_process_time    int64
}

// CrawlProcessor interface
func (h *CrawlHandler)Status(){
}
func (h *CrawlHandler)Process(*proto.CrawlDoc) {
}
func (h *CrawlHandler)Accept(*proto.CrawlDoc) bool{
    return true
}

// CrawlTask Interface
func (h *CrawlHandler)PeriodicTask(p CrawlProcessor){
    for {
        input_chan_size := 0
        if h.input_chan != nil {
            input_chan_size = len(h.input_chan)
        }
        output_chan_size := 0
        if h.output_chan != nil {
            output_chan_size = len(h.output_chan)
        }
        LOG.VLog(3).Debugf("[%s](%d-%d)(%d/%d)(%d/%d)",
                reflect.Indirect(reflect.ValueOf(p)).Type().Name(),
                input_chan_size, output_chan_size,
                h.process_num, h.accept_num,
                h.avg_process_time, h.max_process_time)
        p.Status()
        time_util.Sleep(*CONF.Crawler.PeriodicInterval)
    }
}
func (h *CrawlHandler)Run(p CrawlProcessor) {
    go h.PeriodicTask(p)
    for {
        h.crawlDoc = <- h.input_chan
        h.process_num++
        if p.Accept(h.crawlDoc) {
            now := time_util.GetCurrentTimeStamp()
            p.Process(h.crawlDoc)
            use := time_util.GetCurrentTimeStamp() - now
            if use > h.max_process_time {
                h.max_process_time = use
            }
            h.avg_process_time = ((h.avg_process_time * h.accept_num) + use)/(h.accept_num + 1)
            h.accept_num++
            LOG.VLog(3).Debugf("[%s]Process One Doc %s ",
                reflect.Indirect(reflect.ValueOf(p)).Type().Name(),
                h.crawlDoc.Url)
        }
        h.Output(h.crawlDoc)
    }
}

func (cp *CrawlHandler)Output(doc *proto.CrawlDoc) {
    if cp.output_chan != nil {
        cp.output_chan <- doc
    } else {
        *doc = proto.CrawlDoc{}
    }
}
func (cp *CrawlHandler)SetInputChan(in <-chan*proto.CrawlDoc) {
    cp.input_chan = in
}
func (cp *CrawlHandler)SetOutputChan(out chan<-*proto.CrawlDoc) {
    cp.output_chan = out
}
func (h *CrawlHandler)GetInputChan() <-chan *proto.CrawlDoc{
    return h.input_chan
}
func (h *CrawlHandler)GetOutputChan() chan<- *proto.CrawlDoc{
    return h.output_chan
}
