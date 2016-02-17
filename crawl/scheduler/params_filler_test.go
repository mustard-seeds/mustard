package scheduler
import (
	"testing"
	pb "mustard/crawl/proto"
)
func getCrawlDoc(url string) *pb.CrawlDoc{
	f := PrepareParamFiller{}
	doc := pb.CrawlDoc{RequestUrl:url}
	f.Fill(nil, &doc)
	return &doc
}
func TestPrepareParamFiller(t *testing.T) {
	f := PrepareParamFiller{}
	doc := pb.CrawlDoc{}
	doc.RequestUrl = "http://a.com/"
	f.Fill(nil, &doc)
	if doc.GetCrawlParam() == nil {
		t.Error("PrepareParamFiller Not Fill CrawlParam")
	}
	if doc.GetCrawlRecord() == nil {
		t.Error("PrepareParamFiller Not Fill CrawlRecord")
	}
	if doc.GetCrawlParam() != nil &&doc.CrawlParam.GetFetchHint() == nil {
		t.Error("PrepareParamFiller Not Fill FetchHint")
	}
}
func TestFakeHostParamFiller(t *testing.T) {
	f := FakeHostParamFiller{fakehost:make(map[string]string)}
	f.fakehost["\\w+.sina.com"] = "fake_host.sina.com"
	doc := getCrawlDoc("http://a12.sina.com")
	f.fill(nil,doc)
	if doc.CrawlParam.FakeHost != "fake_host.sina.com" {
		t.Error("FakeHostParamFiller Nof Fill FakeHost")
	}
}
func TestHostLoadParamFiller(t *testing.T) {
	f := HostLoadParamFiller{hostload:make(map[string]int)}
	f.hostload["a.com"] = 1
	doc := getCrawlDoc("http://a.com/")
	f.fill(nil, doc)
	if doc.CrawlParam.Hostload != 1 {
		t.Error("HostLoadParamFiller not Fill Hostload")
	}
	doc1 := getCrawlDoc("http://a.a.com/")
	f.fill(nil, doc1)
	if doc1.CrawlParam.Hostload != int32(*CONF.Crawler.DefaultHostLoad) {
		t.Error("HostLoadParamFiller not Fill Default Hostload")
	}
}
func TestMultiFetcherParamFiller(t *testing.T) {
	f := MultiFetcherParamFiller{multifetcher:make(map[string]int)}
	f.multifetcher["a.com"] = 9
	doc := getCrawlDoc("http://a.com/")
	f.fill(nil, doc)
	if doc.CrawlParam.FetcherCount != 9 {
		t.Error("MultiFetcherParamFiller Not Fill FetcherCount")
	}
}
func TestReceiverParamFiller(t *testing.T) {
	f := ReceiverParamFiller{receivers:make(map[string]*pb.ConnectionInfo)}
	f.receivers["a:3"] = &pb.ConnectionInfo{Host:"a",Port:3}
	f.receivers["b:4"] = &pb.ConnectionInfo{Host:"b",Port:4}
	doc := getCrawlDoc("http://a.com/")
	f.fill(nil, doc)
	if len(doc.CrawlParam.Receivers) != 2 {
		t.Error("ReceiverParamFiller Not Fill Receivers")
	}
}
func TestTagParamFiller(t *testing.T) {
	f := TagParamFiller{}
	doc := getCrawlDoc("http://a.com/")
	f.Fill(&NormalJobD,doc)
	if doc.CrawlParam.Pri != pb.Priority_NORMAL {
		t.Error("TagParamFiller Not Fill Priority")
	}
	if doc.CrawlParam.PrimaryTag != NormalJobD.PrimeTag {
		t.Error("TagParamFiller Not Fill pri tag")
	}
	doc.CrawlParam.PrimaryTag = "XXX"
	f.Fill(&NormalJobD,doc)
	if doc.CrawlParam.PrimaryTag != "XXX" {
		t.Error("TagParamFiller Not Fill pri tag")
	}
}
