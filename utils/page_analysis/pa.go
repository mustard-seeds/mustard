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

/*
  Html Selector
        $("p")  <p> elements.
        $(".test")  all elements with class="test".
        $("#test")  the element with id="test".
*/

type RegexCallback struct {
        regex   string
        callback        func(int, []string)
}
type SelectorCallback struct {
        selector        string
        keyWord         string
        callback        func(int, *goquery.Selection)
}
type HtmlParser struct {
        _content string
        _url *url.URL
        _doc *goquery.Document
        _regCallbacks []*RegexCallback
        _selectorCallbacks []*SelectorCallback
}

func (m *HtmlParser) Reset()  { *m = HtmlParser{} }

// TODO(gaolichuang) registe, parse and callback
// callback will be called multi time for each
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
        // TODO(gaolichuang): optimize to travel DOM and parse each field
        // regex parse
        for _,c := range p._regCallbacks {
                r,_ := regexp.Compile(c.regex)
                regexRet := r.FindAllStringSubmatch(p._content, -1)
                for i,reg := range regexRet {
                        c.callback(i, reg)
                }
        }
        // selector parse
        for _,se := range p._selectorCallbacks {
                p._doc.Find(se.selector).Each(func(i int, s *goquery.Selection){
                        if len(se.keyWord) == 0 {
                                se.callback(i, s)
                        } else {
                                if strings.Contains(s.Text(), se.keyWord) {
                                        se.callback(i, s)
                                }
                        }
                })
        }
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
        p.parseInternal()
        return true, nil
}

func (p *HtmlParser) GetDocument() *goquery.Document {
        return p._doc
}
