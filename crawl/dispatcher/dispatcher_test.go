package dispatcher

import (
	"mustard/crawl/proto"
	"testing"
)

func TestCrawlFeeder(t *testing.T) {
	doc := &proto.CrawlDoc{
		RequestUrl: "http://www.a.com/index.html",
		CrawlParam: &proto.CrawlParam{
			Pri: proto.Priority_NORMAL,
		},
	}
	feeder := CrawlerFeeder{}
	if true != feeder.ShouldOutput(doc) {
		t.Error("Should output check error.")
	}
	*CONF.Crawler.FeederMaxPending = 0
	if false != feeder.ShouldOutput(doc) {
		t.Error("priority normal should not output")
	}
	doc.CrawlParam.Pri = proto.Priority_URGENT
	if true != feeder.ShouldOutput(doc) {
		t.Error("Urgent doc should output.")
	}
	doc.CrawlParam.Pri = proto.Priority_NORMAL
	if false != feeder.AddFeed(doc) || feeder.queuefull_urls != 1 {
		t.Error("Add Feed Check error.")
	}
}
