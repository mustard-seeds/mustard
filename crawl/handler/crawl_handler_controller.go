package handler

import (
    "strings"
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
    "mustard/crawl/handler/handler"
    "reflect"
    "fmt"
)
var CONF = conf.Conf

var InputProcessors []handler.CrawlProcessor
var ProcessChain []handler.CrawlProcessor

func InitCrawlService() {
    for _,name := range strings.Split(*CONF.Crawler.CrawlHandlerChain,";") {
        LOG.Infof("%s Join Crawl Handler Chain", name)
        h := handler.GetCrawlHandlerByName(name).(handler.CrawlHandler)
        p := handler.CrawlHandlerProcessor{}
        p.SetHandler(h)
        ProcessChain = append(ProcessChain, &p)
    }
    if len(ProcessChain) == 0 {
        LOG.Fatal("Crawl handler Chain not assign")
    }
    in := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
    // set input handlers
    for _,name := range strings.Split(*CONF.Crawler.CrawlInputHandler,";") {
        LOG.Infof("%s Input Crawl Processor Start", name)
        r := handler.GetCrawlHandlerByName(name).(handler.CrawlHandlerProcessor)
        fmt.Println(r)
        //if r == nil {
        //    LOG.Fatalf("Can not get struct %s", name)
        //}
        InputProcessors = append(InputProcessors, &r)
        r.SetInputChan(in)
    }
    ProcessChain[0].SetInputChan(in)
    for i := 1;i < len(ProcessChain);i++ {
        out := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
        ProcessChain[i-1].SetOutputChan(out)
        ProcessChain[i].SetInputChan(out)
    }
    ProcessChain[len(ProcessChain)-1].SetOutputChan(make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize))
    for _,p := range ProcessChain {
        LOG.Infof("%s Start to Run", reflect.TypeOf(p.GetHandler()))
        go p.Run()
    }
    for _,r := range InputProcessors {
        go r.Run()
    }
}