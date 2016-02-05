package base

import (
	"mustard/crawl/proto"
	"testing"
)

func TestGetHostName(t *testing.T) {
	doc := &proto.CrawlDoc{
		CrawlParam:&proto.CrawlParam{
			FetchHint:&proto.FetchHint{
				Host:"x.com",
			},
		},
	}
	h := GetHostName(doc)
	if h != "x.com" {
		t.Error("GetHostName without FakeHost not work")
	}
	doc.CrawlParam.FakeHost = "a.com"
	hh := GetHostName(doc)
	if hh != "a.com" {
		t.Error("GetHostName with FakeHost not work")
	}
}
