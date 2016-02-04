package handler

import (
    LOG "mustard/base/log"
    "mustard/crawl/proto"
    "net/http"
    "mustard/base/time_util"
    "mustard/crawl/base"
)

const (
    CONNECTION_POOL_RECOVER_INTERVAL = 3600
    CONNECTION_POOL_TIMEOUT = 3600
)
type Connection struct {
    client  *http.Client
}
// http proxy
// custom
// timeout,connection/read
// follow redirect?

// run in goroutine
func (c *Connection)FetchOne(doc *proto.CrawlDoc, f func(_doc *proto.CrawlDoc, conn *Connection)) {
    // TODO fetch doc and fill field
    time_util.Sleep(3)
    f(doc, c)
}

type ConnectionPool struct {
    record map[string]int64
    hold map[string]bool  // only one host could get in connection pool
    free []*Connection
    busy map[*Connection]bool  // make it could delete
    output_chan  chan<- *proto.CrawlDoc
    last_recover_timestamp int64
}

func (c *ConnectionPool)SetOutChan(output_chan  chan<- *proto.CrawlDoc) {
    c.output_chan = output_chan
}

func (c *ConnectionPool)GetCrawlHostMap() map[string]int64 {
    // TODO release c.record
    return c.record
}
func (c *ConnectionPool)FreeConnectionNum() int {
    return len(c.free)
}
func (c *ConnectionPool)BusyConnectionNum() int {
    return len(c.busy)
}
func (c *ConnectionPool)releaseRecordAndHold() {
    now := time_util.GetCurrentTimeStamp()
    if now - c.last_recover_timestamp < CONNECTION_POOL_RECOVER_INTERVAL {
        return
    }
    release := make([]string, 0)
    for k,v := range c.record {
        if time_util.GetCurrentTimeStamp() - v > CONNECTION_POOL_TIMEOUT {
            release = append(release,k)
        }
    }
    for _,k := range release {
        delete(c.record,k)
        delete(c.hold,k)
        LOG.VLog(3).Debugf("Release Connection Pool Size: %d", len(release))
    }
}
// return false: connection all busy, can not fetch
func (c *ConnectionPool)Fetch(doc *proto.CrawlDoc) bool {
    // check hold or not
    host := base.GetHostName(doc)
    if c.hold[host] == true {
        return false
    }
    if len(c.free) == 0 {
        if len(c.free) + len(c.busy) < *CONF.Crawler.FetchConnectionNum {
            // new dozen conns
            for i := 0;i < 10;i++ {
                conn := &Connection{
                    client: &http.Client{},
                }
                c.free = append(c.free, conn)
            }
        } else {
            LOG.VLog(2).Debugf("Connection Pool full %s/%s",len(c.free), len(c.busy))
            return false
        }
    }
    conn := c.free[0]
    c.free = c.free[1:]
    c.busy[conn] = true
    c.hold[host] = true
    go conn.FetchOne(doc,func(doc *proto.CrawlDoc, conn *Connection){
        c.free = append(c.free, conn)
        delete(c.busy,conn)
        c.record[base.GetHostName(doc)] = time_util.GetCurrentTimeStamp()
        c.hold[base.GetHostName(doc)] = false
        c.output_chan <- doc
    })
    c.releaseRecordAndHold()
    return true
}
