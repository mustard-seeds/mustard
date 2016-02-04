package handler

import (
    "strings"
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
    "mustard/crawl/handler/handler"
    "reflect"
    "mustard/base/string_util"
    "mustard/utils/babysitter"
)

var CONF = conf.Conf

type CrawlHandlerController struct {
     InputProcessors []handler.CrawlTask
     ProcessChain []handler.CrawlTask
}
func (c *CrawlHandlerController)getOutputStat(stat *string, separatorLine, separatorColumn string) {
    *stat = separatorLine + "InputProcessors" + separatorLine
    for _,h := range c.InputProcessors {
        string_util.StringAppendF(stat, "%s:",
            reflect.Indirect(reflect.ValueOf(h)).Type().Name())
        h.Status(stat)
        *stat += separatorColumn
    }
    *stat += separatorLine + "ProcessChain" + separatorLine
    for _,h := range c.ProcessChain {
        string_util.StringAppendF(stat, "%s:",
            reflect.Indirect(reflect.ValueOf(h)).Type().Name())
        h.Status(stat)
        *stat += separatorColumn
    }
}
func (d *CrawlHandlerController)MonitorReportHealthy() error {
    return nil
}
func (c *CrawlHandlerController)MonitorReport(result *babysitter.MonitorResult) {
    // TODO add logic for handler babysitter
    stat := ""
    c.getOutputStat(&stat, "<br>", "<br>")
    result.AddString(stat)
}
func (c *CrawlHandlerController)InitCrawlService() {
    for _,name := range strings.Split(*CONF.Crawler.CrawlHandlerChain,";") {
        LOG.Infof("%s Join Crawl Handler Chain", name)
        h := handler.GetCrawlHandlerByName(name)
        if h == nil {
            LOG.Fatalf("Can not get Crawl Handler %s", name)
        }
        c.ProcessChain = append(c.ProcessChain, h)
    }
    if len(c.ProcessChain) == 0 {
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
        c.InputProcessors = append(c.InputProcessors, r)
        r.SetOutputChan(in)
    }
    c.ProcessChain[0].SetInputChan(in)
    for i := 1;i < len(c.ProcessChain);i++ {
        out := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
        c.ProcessChain[i-1].SetOutputChan(out)
        c.ProcessChain[i].SetInputChan(out)
    }
    c.ProcessChain[len(c.ProcessChain)-1].SetOutputChan(nil)
    for _,p := range c.ProcessChain {
        LOG.Infof("%s Start to Run", reflect.TypeOf(p))
        go p.Run(p.(handler.CrawlProcessor))
    }
    for _,r := range c.InputProcessors {
        LOG.Infof("%s Start to Run", reflect.TypeOf(r))
        go r.Run(r.(handler.CrawlProcessor))
    }
}
func (c *CrawlHandlerController)PrintStatus() {
    stat := ""
    c.getOutputStat(&stat, "\n", "   ")
    LOG.VLog(3).Debug(stat)
}