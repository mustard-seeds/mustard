package page_analysis

import (
        "regexp"
        "strings"
        "net/url"
        "mustard/base/log"
        "mustard/internal/golang.org/x/net/html"
        "mustard/internal/github.com/PuerkitoBio/goquery"
)
var LOG = log.Log

type RegexCallback struct {
        regex   string
        callback        func(int, []string)
}
type SelectorCallback struct {
	selector	string
	callback	func(int, *goquery.Selection)
}
type HtmlParser struct {
        _content string
        _done bool
        _url *url.URL
        _doc *goquery.Document
        _regCallbacks []*RegexCallback
        _selectorCallbacks []*SelectorCallback
}


// TODO(gaolichuang) registe, parse and callback
// the call back cloud be call multi times !!!
func (p *HtmlParser)RegisterRegex(regex string, callback func(int, []string)) {
        p._regCallbacks = append(p._regCallbacks, &RegexCallback{regex, callback})
}

func (p *HtmlParser)RegisterSelectorWithTextKeyWord(selector string, keyword string, callback func(int, *goquery.Selection)) {
	p._selectorCallbacks = append(p._selectorCallbacks, &SelectorCallback{selector, keyword, callback})
}
func (p *HtmlParser)RegisterSelector(selector string, callback func(int, *goquery.Selection)) {
	p.RegisterSelectorWithTextKeyWord(selector, "", callback)
}

func (p *HtmlParser)parseInternal() {

}
func (p *HtmlParser)Parse(_url, _content string) (result bool, err error) {
        var u *url.URL
        var e error
        if u, e = url.Parse(_url);e != nil {
                LOG.Error("Parse Url Fail")
                return false, e
        }
        p._url = u
        p._content = _content
        var root *html.Node
        // strings.NewReader, make string like io.Reader
        if root, e = html.Parse(strings.NewReader(p._content)); e != nil {
                LOG.Error("parse fail. url:" + p._url.String())
                return false, e
        }
        p._doc = goquery.NewDocumentFromNode(root)
        for _,c := range p._regCallbacks {
		        c.callback("hello")
        }
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
