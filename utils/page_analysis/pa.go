package page_analysis

import (
	"strings"
	"net/url"
	"mustard/base/log"
	"mustard/internal/golang.org/x/net/html"
	"mustard/internal/github.com/PuerkitoBio/goquery"
)
var LOG = log.Log

type HtmlParser struct {
	_content string
	_done bool
	_url *url.URL
	_doc *goquery.Document
}


func (p *HtmlParser) Parse(_url, _content string) (result bool, err error) {
	// parse url
	var u *url.URL
	var e error
	if u, e = url.Parse(_url);e != nil {
		LOG.Error("Parse Url Fail")
		return false, e
	}
	p._url = u
	p._content = _content
	var root *html.Node
	if root, e = html.Parse(strings.NewReader(p._content)); e != nil {
		LOG.Error("parse fail. url:" + p._url.String())
		return false, e
	}
	p._doc = goquery.NewDocumentFromNode(root)
	return true, nil
}
/*
$("p")  <p> elements.
$(".test")  all elements with class="test".
$("#test")  the element with id="test".
*/
func (p *HtmlParser) GetLinkByHost(host string) map[string]string {
	var links = map[string]string{}
	p._doc.Find("a").Each(func(i int, s *goquery.Selection){
		href , _ := s.Attr("href")
		if strings.HasPrefix(href,"/") {
			if host == "" || strings.Compare(p._url.Host, host) == 0 {
				links["http://" + p._url.Host + href] = s.Text()
			}
		} else if strings.HasPrefix(href, "http") {
			if u, e := url.Parse(href);e != nil {
				return
			} else {
				if host == "" || strings.Compare(u.Host, host) == 0 {
					links[u.String()] = s.Text()
				}
			}
		}
	})
	return links
}
func (p *HtmlParser) GetDocument() *goquery.Document {
	return p._doc
}
