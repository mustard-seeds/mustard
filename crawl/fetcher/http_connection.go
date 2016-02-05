package fetcher

import (
    pb "mustard/crawl/proto"
    "net/http"
    "mustard/base/time_util"
)

type Connection struct {
    client  *http.Client
}
// use proxy or not.
// custom
// timeout,connection/read
// follow redirect?
// basic auth
// cookie

// run in goroutine
func (c *Connection)FetchOne(doc *pb.CrawlDoc, f func(*pb.CrawlDoc, *Connection)) {
    // TODO fetch doc and fill field
    time_util.Sleep(3)
    f(doc, c)
}

func NewConnection() *Connection {
    return &Connection{
        client: &http.Client{},
    }
}