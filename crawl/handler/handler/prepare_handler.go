package handler

import (
	LOG "mustard/base/log"
	"mustard/crawl/proto"
	"mustard/crawl/base"
	"mustard/internal/golang.org/x/text/transform"
	"mustard/internal/golang.org/x/net/html/charset"
	"strings"
	"io/ioutil"
)

type PrepareHandler struct {
	CrawlHandler
}

func (handler *PrepareHandler)Accept(crawlDoc *proto.CrawlDoc) bool {
	return true
}

func (handler *PrepareHandler)Process(crawlDoc *proto.CrawlDoc) {
	LOG.VLog(4).Debugf("\n%s", base.DumpCrawlDoc(crawlDoc))
	// dump unicode content..
	LOG.VLog(4).Debugf("Content:\n%U", crawlDoc.Content)
	crawlDoc.ContentLength = len(crawlDoc.Content)
	// charset detect if not utf-8, decode it.
	translateEncoding(crawlDoc)

}
func translateEncoding(crawlDoc *proto.CrawlDoc) {
	e, n, _ := charset.DetermineEncoding(crawlDoc.Content, crawlDoc.ContentType)
	crawlDoc.OrigEncoding = n
	if n != "utf-8" {
		if e == nil {
			crawlDoc.ConvEncoding = n
		} else {
			s, err := transformString(e.NewDecoder(), crawlDoc.Content)
			if err != nil {
				crawlDoc.ConvEncoding = n
			} else {
				crawlDoc.ConvEncoding = "utf-8"
				crawlDoc.Content = s
			}
		}
	} else {
		crawlDoc.ConvEncoding = n
	}
}
func transformString(t transform.Transformer, s string) (string, error) {
	r := transform.NewReader(strings.NewReader(s), t)
	b, err := ioutil.ReadAll(r)
	return string(b), err
}

// use for create instance from a string
func init() {
	registerCrawlTaskType(&PrepareHandler{})
}
