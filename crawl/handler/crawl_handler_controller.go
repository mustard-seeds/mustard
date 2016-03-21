package handler

import (
	"mustard/base"
	"mustard/base/conf"
	LOG "mustard/base/log"
	"mustard/base/string_util"
	"mustard/crawl/handler/handler"
	"mustard/crawl/proto"
	"mustard/utils/babysitter"
	"reflect"
	"strings"
)

var CONF = conf.Conf

type CrawlHandlerController struct {
	InputProcessors []handler.CrawlTask
	ProcessChain    []handler.CrawlTask
	inited          bool
}

func (c *CrawlHandlerController) getOutputStat(stat *string, separatorLine, separatorColumn string) {
	*stat = separatorLine + "InputProcessors" + separatorLine
	for _, h := range c.InputProcessors {
		string_util.StringAppendF(stat, "%s:",
			reflect.Indirect(reflect.ValueOf(h)).Type().Name())
		h.Status(stat)
		*stat += separatorColumn
	}
	*stat += separatorLine + "ProcessChain" + separatorLine
	for _, h := range c.ProcessChain {
		string_util.StringAppendF(stat, "%s:",
			reflect.Indirect(reflect.ValueOf(h)).Type().Name())
		h.Status(stat)
		*stat += separatorColumn
	}
}
func (d *CrawlHandlerController) MonitorReportHealthy() error {
	// TODO: handler healthy defination...
	return nil
}
func (c *CrawlHandlerController) MonitorReport(result *babysitter.MonitorResult) {
	// TODO add logic for handler babysitter
	stat := ""
	c.getOutputStat(&stat, "<br>", "<br>")
	result.AddString(stat)
}
func (c *CrawlHandlerController) initServiceInternal() {
	for _, name := range strings.Split(*CONF.Crawler.CrawlHandlerChain, ";") {
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
	for _, name := range strings.Split(*CONF.Crawler.CrawlInputHandler, ";") {
		LOG.Infof("%s Input Crawl Processor Start", name)
		r := handler.GetCrawlHandlerByName(name)
		if r == nil {
			LOG.Fatalf("Can not get crawl processor %s", name)
		}
		c.InputProcessors = append(c.InputProcessors, r)
		r.SetOutputChan(in)
	}
	c.ProcessChain[0].SetInputChan(in)
	for i := 1; i < len(c.ProcessChain); i++ {
		out := make(chan *proto.CrawlDoc, *CONF.Crawler.ChannelBufSize)
		c.ProcessChain[i-1].SetOutputChan(out)
		c.ProcessChain[i].SetInputChan(out)
	}
	c.ProcessChain[len(c.ProcessChain)-1].SetOutputChan(nil)
}
func (c *CrawlHandlerController) startServiceInterval() {
	for _, p := range c.ProcessChain {
		LOG.Infof("%s Start to Run", reflect.TypeOf(p))
		base.CHECK(p.Init(), "%s Init Fail", reflect.TypeOf(p))
		go p.Run(p.(handler.CrawlProcessor))
	}
	for _, r := range c.InputProcessors {
		LOG.Infof("%s Start to Run", reflect.TypeOf(r))
		base.CHECK(r.Init(), "%s Init Fail", reflect.TypeOf(r))
		go r.Run(r.(handler.CrawlProcessor))
	}
}
func (c *CrawlHandlerController) InitCrawlService() {
	if c.inited == true {
		base.CHECK(false, "Handler Controller should only Init Once.")
	}
	c.inited = true
	c.initServiceInternal()
	c.startServiceInterval()
}
func (c *CrawlHandlerController) PrintStatus() {
	stat := ""
	c.getOutputStat(&stat, "\n", "   ")
	LOG.VLog(2).Debug(stat)
}
