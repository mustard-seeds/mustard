package fetcher

import (
	"mustard/crawl/proto"
	"testing"
)

func TestHostLoadQueue(t *testing.T) {
	hlq := newHostLoadQueue()
	doc := &proto.CrawlDoc{RequestUrl:"http://a.com/"}
	if !hlq.Empty() {
		t.Error("HostLoadQueue not empty after init?")
	}
	hlq.Push(doc)
	if 1 != hlq.Size() {
		t.Error("HostLoadQueue Size not right.")
	}
	top,err := hlq.Top()
	if err != nil {
		t.Error("Top get Error.")
	}
	if top.RequestUrl != doc.RequestUrl {
		t.Error("Top element error")
	}
	hlq.Pop()
	if !hlq.Empty() {
		t.Error("Pop Error")
	}
	hlq.capacity = 1
	hlq.Push(doc)
	hlq.Push(doc)
	err = hlq.Push(doc)
	if err == nil {
		t.Error("Push error, it should be full.")
	}
}
