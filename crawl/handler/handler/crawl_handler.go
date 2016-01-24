package handler

import (
    "time"
    "reflect"
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
)

var CONF = conf.Conf

type CrawlHandler interface {
    Accept(*proto.CrawlDoc) bool
    Process(*proto.CrawlDoc)
    Status()
}
type CrawlProcessor interface {
    GetInputChan() <-chan *proto.CrawlDoc
    GetOutputChan() chan<- *proto.CrawlDoc
    SetInputChan(in <-chan*proto.CrawlDoc)
    SetOutputChan(out chan<-*proto.CrawlDoc)
    Output(doc *proto.CrawlDoc)
    GetHandler() CrawlHandler
    SetHandler(CrawlHandler)
    Run()
}

type CrawlHandlerProcessor struct {
    handler CrawlHandler
    input_chan <-chan *proto.CrawlDoc
    output_chan  chan<- *proto.CrawlDoc
    crawlDoc *proto.CrawlDoc
    process_num int
    accept_num  int
}
func (cp *CrawlHandlerProcessor)GetInputChan() <-chan *proto.CrawlDoc{
    return cp.input_chan
}
func (cp *CrawlHandlerProcessor)GetOutputChan() chan<- *proto.CrawlDoc{
    return cp.output_chan
}
func (cp *CrawlHandlerProcessor)GetHandler() CrawlHandler{
    return cp.handler
}
func (cp *CrawlHandlerProcessor)PeriodicTask(){
    for {
        input_chan_size := 0
        if cp.input_chan != nil {
            input_chan_size = len(cp.input_chan)
        }
        output_chan_size := 0
        if cp.output_chan != nil {
            output_chan_size = len(cp.output_chan)
        }
        LOG.VLog(3).Debugf("[%s]in(%d)-out(%d)",reflect.TypeOf(cp.handler),input_chan_size,output_chan_size)
        cp.handler.Status()
        time.Sleep(time.Second * time.Duration(*CONF.Crawler.PeriodicInterval))
    }
}
func (cp *CrawlHandlerProcessor)Run() {
    if cp.handler == nil {
        LOG.Fatal("Crawler Should assign One Handler")
    }
    go cp.PeriodicTask()
    for {
        cp.crawlDoc = <- cp.input_chan
        cp.process_num++
        if cp.handler.Accept(cp.crawlDoc) {
            cp.accept_num++
            LOG.VLog(3).Debugf("Process One Doc %s ",cp.crawlDoc.Url)
            cp.handler.Process(cp.crawlDoc)
        }
        cp.Output(cp.crawlDoc)
    }
}
func (cp *CrawlHandlerProcessor)SetHandler(handler CrawlHandler) {
    cp.handler = handler
}
func (cp *CrawlHandlerProcessor)SetInputChan(in <-chan*proto.CrawlDoc) {
    cp.input_chan   = in
}
func (cp *CrawlHandlerProcessor)SetOutputChan(out chan<-*proto.CrawlDoc) {
    cp.output_chan = out
}
func (cp *CrawlHandlerProcessor)Output(doc *proto.CrawlDoc) {
    if cp.output_chan != nil {
        cp.output_chan <- doc
    }
}
