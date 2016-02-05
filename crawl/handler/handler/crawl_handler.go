package handler

import (
    "mustard/base/time_util"
    "mustard/base/conf"
    "reflect"
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "mustard/base/string_util"
    "fmt"
)

var CONF = conf.Conf

type CrawlTask interface {
    Init() bool             // call at crawlhandlercontroller before Run
    Run(p CrawlProcessor)   // you can overwrite run function.
    Output(doc *proto.CrawlDoc)
    SetInputChan(in <-chan*proto.CrawlDoc)
    SetOutputChan(out chan<-*proto.CrawlDoc)
    GetInputChan() <-chan *proto.CrawlDoc
    GetOutputChan() chan<- *proto.CrawlDoc
    Status(s *string)
}

type CrawlProcessor interface {
    Process(*proto.CrawlDoc)  // use CrawlHandler Run Func
    Accept(*proto.CrawlDoc) bool // use Crawlhandler Run Func
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
func (h *CrawlHandler)Init() bool {
    return true
}
// CrawlProcessor interface
func (h *CrawlHandler)Status(s *string){
    ins := "X"
    if h.input_chan != nil {
        ins = fmt.Sprintf("%d",len(h.input_chan))
    }
    string_util.StringAppendF(s, "(%s)(%d/%d %d/%d)", ins,
        h.process_num, h.accept_num,
        h.avg_process_time, h.max_process_time)
}

func (h *CrawlHandler)Process(*proto.CrawlDoc) {
}
func (h *CrawlHandler)Accept(*proto.CrawlDoc) bool{
    return true
}

func (h *CrawlHandler)Run(p CrawlProcessor) {
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
