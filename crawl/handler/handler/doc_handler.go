/*
Page Analysis and extract link.
*/
package handler

import (
	"fmt"
	LOG "mustard/base/log"
	"mustard/base/string_util"
	"mustard/crawl/base"
	"mustard/crawl/proto"
	"mustard/internal/github.com/PuerkitoBio/goquery"
	"mustard/utils/page_analysis"
	"mustard/utils/url_parser"
	"net/url"
	"reflect"
	"strings"
)

type DocHandler struct {
	CrawlHandler
	htmlParser *page_analysis.HtmlParser
	doc        *proto.CrawlDoc
	domain     string
}

func (handler *DocHandler) Init() bool {
	handler.htmlParser = page_analysis.NewHtmlParser()
	handler.htmlParser.RegisterSelector("a", handler.extractLinkCallBack)
	return true
}
func (handler *DocHandler) extractLinkCallBack(i int, s *goquery.Selection) {
	href, hrefexit := s.Attr("href")
	if !hrefexit {
		return
	}
	if !(strings.HasPrefix(href, "/") || strings.HasPrefix(href, "http") || strings.HasPrefix(href, ".")) {
		LOG.VLog(4).Debugf("Not Avaliable link %s", href)
		return
	}
	nofollow, nofollowexit := s.Attr("rel")
	if (!handler.doc.CrawlParam.Nofollow) && nofollowexit && nofollow == "nofollow" {
		LOG.VLog(4).Debugf("NoFollow link doc.crawlparam.nofollow:%t, link:%s,text:%s",
			handler.doc.CrawlParam.Nofollow,
			href,
			s.Text())
		return
	}
	text := string_util.Purify(s.Text(), "\n", "\t", " ")
	if strings.HasPrefix(href, "/") {
		requrl, _ := url.Parse(handler.doc.RequestUrl)
		LOG.VLog(4).Debugf("InDomainLinkFill %s,text:%U", href, text)
		handler.doc.IndomainOutlinks = append(handler.doc.IndomainOutlinks, &proto.OutLink{
			Url:  fmt.Sprintf("%s://%s%s", requrl.Scheme, requrl.Host, href),
			Text: text,
		})
	} else {
		newdomain := base.GetDomainFromHost(url_parser.GetHost(href))
		if base.IsSameDomain(newdomain, handler.domain) {
			LOG.VLog(4).Debugf("InDomainLink %s,text:%U", href, text)
			handler.doc.IndomainOutlinks = append(handler.doc.IndomainOutlinks, &proto.OutLink{
				Url:  href,
				Text: text,
			})
		} else {
			LOG.VLog(4).Debugf("OutDomainLink %s,text:%U", href, text)
			handler.doc.OutdomainOutlinks = append(handler.doc.OutdomainOutlinks, &proto.OutLink{
				Url:  href,
				Text: text,
			})
		}
	}
}
func (handler *DocHandler) Accept(crawlDoc *proto.CrawlDoc) bool {
	return base.IsCrawlSuccess(crawlDoc.Code)
}
func (handler *DocHandler) Process(crawlDoc *proto.CrawlDoc) {
	LOG.VLog(3).Debugf("[%s]Process One Doc %s ",
		reflect.Indirect(reflect.ValueOf(handler)).Type().Name(),
		crawlDoc.Url)
	LOG.VLog(4).Debugf("DocHandler. DumpCrawlDoc\n%s", base.DumpCrawlDoc(crawlDoc))
	handler.doc = crawlDoc
	handler.domain = base.GetDomainFromHost(url_parser.GetHost(crawlDoc.RequestUrl))

	handler.htmlParser.Parse(handler.doc.Url, handler.doc.Content)
}

// use for create instance from a string
func init() {
	registerCrawlTaskType(&DocHandler{})
}
