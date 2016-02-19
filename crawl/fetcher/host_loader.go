package fetcher

import (
    LOG "mustard/base/log"
    "mustard/crawl/base"
    "mustard/crawl/proto"
    "errors"
    "mustard/base/time_util"
    "mustard/base/conf"
    "mustard/base/hash"
)
var CONF = conf.Conf

/*
    hostload and connection pool
*/
type HostLoadQueue struct {
    hosts []*proto.CrawlDoc
    capacity    int
}
func (hlq *HostLoadQueue)Top() (*proto.CrawlDoc,error) {
    if hlq.Empty() {
        return nil,errors.New("Host Queue Empty")
    }
    doc := hlq.hosts[0]
    return doc,nil
}
func (hlq *HostLoadQueue)Push(doc *proto.CrawlDoc) error {
    if hlq.Full() {
        return errors.New("Host Queue Full")
    }
    hlq.hosts = append(hlq.hosts, doc)
    return nil
}

func (hlq *HostLoadQueue)Pop() (*proto.CrawlDoc,error) {
    t,e := hlq.Top()
    if e != nil {
        return t,e
    }
    hlq.hosts =hlq.hosts[1:]
    return t,e
}
func (hlq *HostLoadQueue)Full() bool {
    return hlq.Size() > hlq.capacity
}
func (hlq *HostLoadQueue)Empty() bool {
    return 0 == hlq.Size()
}
func (hlq *HostLoadQueue)Size() int {
    return len(hlq.hosts)
}

//////////////////////////////////////////////////////////
// no need lock, it is all in one thread.
type HostLoader struct {
    hostMap map[string]*HostLoadQueue
    uim int
    him int
}
func (hl *HostLoader)Uim() int {
    return hl.uim
}
func (hl *HostLoader)Him() int {
    return hl.him
}
func (hl *HostLoader)Push(doc *proto.CrawlDoc) error {
    host := base.GetHostName(doc)
    q,exist := hl.hostMap[host]
    if exist {
        if q.Full() {
            return errors.New(host + "QueueFull")
        } else {
            q.Push(doc)
        }
    } else {
        hl.hostMap[host] = newHostLoadQueue()
        hl.hostMap[host].Push(doc)
    }
    return nil
}
// param s: host last crawl time
// param f: callback for already reach time crawldoc
func (hl *HostLoader)Travel(s map[string]int64, f func(*proto.CrawlDoc) bool) {
    // release hostMap
    // call func f for avaliable doc
    var rel,lens,docs int
    release := make([]string, 0)
    for k,v := range hl.hostMap {
        lens++
        docs += v.Size()
        if v.Empty() {
            rel++
            release = append(release, k)
            continue
        }
        now := time_util.GetCurrentTimeStamp()
        doc,_ := v.Top()
        if now - s[base.GetHostName(doc)] > int64(int(doc.CrawlParam.Hostload) + hash.RandomIntn(1 + int(doc.CrawlParam.RandomHostload))) {
            if f(doc) {
                v.Pop()
            }
        }
    }
    hl.him = lens
    hl.uim = docs
    // release or not
    if float64(rel)/float64(lens) > *CONF.Crawler.HostLoaderReleaseRatio {
        for _,k := range release {
            delete(hl.hostMap, k)
        }
        LOG.VLog(3).Debugf("HostLoad release.%s",release)
    }
}

func NewHostLoader() *HostLoader {
    return &HostLoader{
        hostMap: make(map[string]*HostLoadQueue),
    }
}
func newHostLoadQueue() *HostLoadQueue {
    return &HostLoadQueue{
        capacity:*CONF.Crawler.HostLoaderQueueSize,
    }
}