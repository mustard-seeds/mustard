package handler

import (
    "strings"
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
    "mustard/crawl/handler/handler"
    "reflect"
)
var CONF = conf.Conf

var InputProcessors []handler.CrawlTask
var ProcessChain []handler.CrawlTask

func InitCrawlService() {
    for _,name := range strings.Split(*CONF.Crawler.CrawlHandlerChain,";") {
        LOG.Infof("%s Join Crawl Handler Chain", name)
        h := handler.GetCrawlHandlerByName(name)
        if h == nil {
            LOG.Fatalf("Can not get Crawl Handler %s", name)
        }
        ProcessChain = append(ProcessChain, h)
    }
    if len(ProcessChain) == 0 {
        LOG.Fatal("Crawl handler Chain not assign")
    }
    in := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
    // set input handlers
    for _,name := range strings.Split(*CONF.Crawler.CrawlInputHandler,";") {
        LOG.Infof("%s Input Crawl Processor Start", name)
        r := handler.GetCrawlHandlerByName(name)
        if r == nil {
            LOG.Fatalf("Can not get crawl processor %s", name)
        }
        InputProcessors = append(InputProcessors, r)
        r.SetOutputChan(in)
    }
    ProcessChain[0].SetInputChan(in)
    for i := 1;i < len(ProcessChain);i++ {
        out := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
        ProcessChain[i-1].SetOutputChan(out)
        ProcessChain[i].SetInputChan(out)
    }
    ProcessChain[len(ProcessChain)-1].SetOutputChan(nil)
    for _,p := range ProcessChain {
        LOG.Infof("%s Start to Run", reflect.TypeOf(p))
        go p.Run(p.(handler.CrawlProcessor))
    }
    for _,r := range InputProcessors {
        LOG.Infof("%s Start to Run", reflect.TypeOf(r))
        go r.Run(r.(handler.CrawlProcessor))
    }
}