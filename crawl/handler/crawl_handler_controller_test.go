package handler

import (
	"reflect"
	"testing"
	"fmt"
)

func TestInitCrawlService(t *testing.T) {
	*CONF.Crawler.CrawlHandlerChain="FetchHandler;DocHandler;StorageHandler"
	*CONF.Crawler.CrawlInputHandler="RequestProcessor;DummyRequestProcessor"
	c := &CrawlHandlerController{}
	c.initServiceInternal()
	if len(c.ProcessChain) != 3 {
		t.Error("CrawlHandlerInit Fail for Processor Chain")
	}
	if len(c.InputProcessors) != 2 {
		t.Error("CrawlHandlerInit Fail For InputChain")
	}
	fmt.Println(reflect.TypeOf(c.ProcessChain[0]).String())
	if reflect.TypeOf(c.ProcessChain[0]).String() != "*handler.FetchHandler" {

		t.Error("CrawlHandlerInit Sequence error.")
	}
}
