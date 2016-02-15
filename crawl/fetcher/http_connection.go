package fetcher

import (
    pb "mustard/crawl/proto"
    "net/http"
    "mustard/base/time_util"
    "time"
    "net"
    "errors"
)
const (
    BROWSER_UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
    CONNECTION_TIMEOUT time.Duration = time.Duration(3) * time.Second
    READ_TIMEOUT       time.Duration = time.Duration(30) * time.Second
)

var GeneralHeader = map[string]string {
    "Accept":"text/html;q=0.8, */*;q=0.5",
    "Accept-Charset":"utf-8, gbk, gb2312, *;q=0.5",
    "Accept-Language":"zh-cn;q=0.8, *;q=0.5",
    "Accept-Encoding":"gzip",
    "Connection":"close",
//    "Connection":"keep-alive",
	"User-Agent":"MustardSpider",
}
// Referer:
// Authorization: Basic
// Cookie:

type FetchTimeout struct {
    connect     time.Duration
    readwrite   time.Duration
}


type Connection struct {
    client  *http.Client
}
// use proxy or not.
// custom
// timeout,connection/read
// follow redirect?
// basic auth
// cookie
// encode...
// header...
// proxy
// user agent
// Vlog4 print debug info -- format...
// support https.

// CODE -- match info
//  connect timeout -- NOCONNECTION
//  readwrite timeout -- TIMEOUT
//  302 redirect no url or no header. -- BADHEADER

// http://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
func noRedirect(req *http.Request, via []*http.Request) error {
    return errors.New("Don't redirect!")
}

func (c *Connection)TimeoutDialer(to *FetchTimeout) func(net, addr string) (c net.Conn, err error) {
    return nil
}
func (c *Connection)noRedirect(req *http.Request, via []*http.Request) error {
    return errors.New("No Redirect")
}
// run in goroutine
func (c *Connection)FetchOne(doc *pb.CrawlDoc, f func(*pb.CrawlDoc, *Connection)) {
    // TODO fetch doc and fill field
    time_util.Sleep(3)
    f(doc, c)
}

func (c *Connection)Handle200(resp *http.Response, doc *pb.CrawlDoc) {

}
func (c *Connection)Handle30X(resp *http.Response, doc *pb.CrawlDoc) {

}

func NewConnection() *Connection {
    return &Connection{
        client: &http.Client{},
    }
}
