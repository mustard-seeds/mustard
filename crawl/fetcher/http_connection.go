package fetcher

import (
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	LOG "mustard/base/log"
	crawl_base "mustard/crawl/base"
	pb "mustard/crawl/proto"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	BROWSER_UA                         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
	CONNECTION_TIMEOUT   time.Duration = time.Duration(3) * time.Second
	READ_WRITE_TIMEOUT   time.Duration = time.Duration(30) * time.Second
	CHECK_REDIRECT_DEPTH               = 5
)

var GeneralHeader = map[string]string{
	"Accept":          "text/html;q=0.8, */*;q=0.5",
	"Accept-Charset":  "utf-8, gbk, gb2312, *;q=0.5",
	"Accept-Language": "zh-cn;q=0.8, *;q=0.5",
	"Accept-Encoding": "gzip",
	"Connection":      "close",
	//    "Connection":"keep-alive",
	"User-Agent": "MustardSpider",
}

type FetchTimeout struct {
	connect   time.Duration
	readwrite time.Duration
}

var GeneralFetchTime = &FetchTimeout{
	connect:   CONNECTION_TIMEOUT,
	readwrite: READ_WRITE_TIMEOUT,
}

func timeoutDialer(to *FetchTimeout) func(net, addr string) (c net.Conn, err error) {
	return nil
}

// http://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("No Redirect")
}
func multiCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= CHECK_REDIRECT_DEPTH {
		return errors.New(fmt.Sprintf("stopped after %d redirects", CHECK_REDIRECT_DEPTH))
	}
	return nil
}

type Connection struct {
	clientGenerator  *HttpClientGenerator
	requestGenerator *HttpRequestGenerator
}

//  302 redirect no url or no header. -- BADHEADER

// run in goroutine
func (c *Connection) FetchOne(doc *pb.CrawlDoc, f func(*pb.CrawlDoc, *Connection)) {
	// TODO fetch doc and fill field
	// step 1. fill request info
	client := c.clientGenerator.
		WithSchema(doc.RequestUrl).
		WithRedirect(doc.CrawlParam.FollowRedirect).
		WithProxy(doc.CrawlParam.UseProxy).
		NewClient()
	req := c.requestGenerator.
		WithCustomUA(doc.CrawlParam.CustomUa).
		WithReferer(doc.CrawlParam.Referer).
		NewRequest(doc.Url)
	// step 2. fetch
	resp, err := client.Do(req)
	// step 3. judge response code and fill the nessary field of crawldoc
	if resp != nil {
		dumpResp, dumpErr := httputil.DumpResponse(resp, false)
		respMsg, respErr := "Nil", "Nil"
		if dumpResp != nil {
			respMsg = string(dumpResp)
		}
		if dumpErr != nil {
			respErr = dumpErr.Error()
		}
		LOG.VLog(4).Debugf("Dump Response(Error:%s):\n%s", respMsg, respErr)
	}
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
		c.clientGenerator.MarkDeadProxy()
		doc.Code = pb.ReturnType_NOCONNECTION
		c.HandleOther(resp, err, doc)
	} else if err != nil && strings.Contains(err.Error(), "i/o timeout") {
		c.clientGenerator.MarkDeadProxy()
		// read tcp 172.24.47.104:54386->220.181.112.244:443: i/o timeout
		// dial tcp: i/o timeout
		doc.Code = pb.ReturnType_TIMEOUT
		c.HandleOther(resp, err, doc)
	} else if err != nil && strings.Contains(err.Error(), "No Redirect") {
		// redirect error throw.
		doc.Code = pb.ReturnType(resp.StatusCode)
		c.Handle30X(resp, doc)
	} else if err == nil {
		doc.Code = pb.ReturnType(resp.StatusCode)
		if crawl_base.IsCrawlSuccess(pb.ReturnType(resp.StatusCode)) {
			c.Handle200(resp, doc)
		} else {
			c.HandleOther(resp, nil, doc)
		}
	} else {
		c.clientGenerator.MarkDeadProxy()
		// other?
		c.HandleOther(resp, err, doc)
	}
	// the last step: call the callback function.
	f(doc, c)
}

func (c *Connection) Handle200(resp *http.Response, doc *pb.CrawlDoc) {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	if b, err := ioutil.ReadAll(reader); err == nil {
		doc.Content = string(b)
	}

	dumResp, _ := httputil.DumpResponse(resp, false)
	doc.Header = string(dumResp)
	doc.LastModify = resp.Header.Get("last-modified")
	doc.ContentType = resp.Header.Get("Content-Type")
	LOG.VLog(3).Debugf("Fetch Success, url:%s,reqtype:%d", doc.Url, doc.CrawlParam.Rtype)
}

func (c *Connection) Handle30X(resp *http.Response, doc *pb.CrawlDoc) {
	doc.LastModify = resp.Header.Get("last-modified")
	doc.ContentType = resp.Header.Get("Content-Type")
	redirectUrl := resp.Header.Get("Location")
	if !crawl_base.IsInvalidUrl(redirectUrl) {
		doc.Code = pb.ReturnType_INVALIDREDIRECTURL
	} else {
		doc.RedirectUrl = redirectUrl
	}
	LOG.VLog(3).Debugf("Fetch 30X, url:%s, redirecturl:%s, reqtype:%d", doc.Url, doc.RedirectUrl, doc.CrawlParam.Rtype)
}
func (c *Connection) HandleOther(resp *http.Response, err error, doc *pb.CrawlDoc) {
	if err != nil {
		doc.ErrorInfo = err.Error()
	}
	LOG.VLog(3).Debugf("Fetch Code:%d, url:%s, reqtype:%d", doc.Code, doc.Url, doc.CrawlParam.Rtype)
}
func NewConnection() *Connection {
	return &Connection{
		clientGenerator: &HttpClientGenerator{
			httpProxy: NewProxyManager(PROXY_SELECT_RR),
			redirect:  false,
			https:     false,
			proxy:     false,
		},
		requestGenerator: &HttpRequestGenerator{
			customUA: false,
			referer:  "",
		},
	}
}

////////////////HttpClientGenerator//////////////////////////////////////////////////////////
type HttpClientGenerator struct {
	httpProxy *ProxyManager
	redirect  bool
	https     bool // if use https, no proxy...
	proxy     bool
	proxyUrl  *url.URL
}

func (hg *HttpClientGenerator) reset() {
	hg.redirect = false
	hg.https = false
	hg.proxy = false
}
func (hg *HttpClientGenerator) WithSchema(_url string) *HttpClientGenerator {
	if strings.HasPrefix(_url, "https") {
		hg.https = true
	} else {
		hg.https = false
	}
	return hg
}
func (hg *HttpClientGenerator) WithRedirect(y bool) *HttpClientGenerator {
	hg.redirect = y
	return hg
}
func (hg *HttpClientGenerator) WithProxy(y bool) *HttpClientGenerator {
	hg.proxy = y
	return hg
}

func (hg *HttpClientGenerator) NewClient() *http.Client {
	// TODO. add cache for new http client.
	LOG.VLog(4).Debugf("NewClient:https:%t,proxy:%t,redirect:%t", hg.https, hg.proxy, hg.redirect)
	var client *http.Client
	ckRedirect := noRedirect
	if hg.redirect == true {
		ckRedirect = multiCheckRedirect
	}
	hg.proxyUrl = nil

	var tlsClientConfig *tls.Config = nil
	if hg.https {
		tlsClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var clientProxy func(*http.Request) (*url.URL, error) = nil
	if hg.proxy && hg.https == false { // only http request use proxy
		proxyUrl, err := hg.httpProxy.GetProxyUrl()
		if err != nil {
			hg.proxyUrl = proxyUrl
			clientProxy = http.ProxyURL(proxyUrl)
		}
	}
	client = &http.Client{
		CheckRedirect: ckRedirect,
		Transport: &http.Transport{
			Dial:            timeoutDialer(GeneralFetchTime),
			Proxy:           clientProxy,
			TLSClientConfig: tlsClientConfig,
		},
	}
	hg.reset()
	return client
}
func (hg *HttpClientGenerator) MarkDeadProxy() {
	if hg.proxyUrl != nil {
		hg.httpProxy.MarkDeadProxy(hg.proxyUrl)
	}
}

/////////////HttpRequestGenerator/////////////////////////////////////////////////////////////
// HttpRequestGenerator TODO. Add cookie and basic auth support...
type HttpRequestGenerator struct {
	customUA bool
	referer  string
}

func (rg *HttpRequestGenerator) WithCustomUA(y bool) *HttpRequestGenerator {
	rg.customUA = y
	return rg
}
func (rg *HttpRequestGenerator) WithReferer(referer string) *HttpRequestGenerator {
	rg.referer = referer
	return rg
}
func (rg *HttpRequestGenerator) NewRequest(_url string) *http.Request {
	// TODO. support POST method.... FetchHint.post_data
	req, _ := http.NewRequest("GET", _url, nil)

	for k, v := range GeneralHeader {
		req.Header.Set(k, v)
	}
	if rg.referer != "" {
		req.Header.Set("Referer", rg.referer)
	}
	if rg.customUA {
		req.Header.Set("User-Agent", BROWSER_UA)
	}
	dumpReq, _ := httputil.DumpRequest(req, true)
	LOG.VLog(4).Debugf("DumpRequest:\n%s", string(dumpReq))
	return req
}

//////////////////////////////////////////////////////////////////////////
