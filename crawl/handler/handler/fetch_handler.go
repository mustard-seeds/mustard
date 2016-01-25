package handler

import (
    LOG "mustard/base/log"
    "mustard/base/conf"
    "mustard/crawl/proto"
)
var CONF = conf.Conf

type HostLoadQueue struct {
    hosts []*proto.CrawlDoc
}
func (hlq *HostLoadQueue)Push(doc *proto.CrawlDoc) {
    hlq.hosts = append(hlq.hosts, doc)
}
func (hlq *HostLoadQueue)Pop() *proto.CrawlDoc {
    if len(hlq.hosts) == 0 {
        return nil
    }
    doc := hlq.hosts[0]
    hlq.hosts =hlq.hosts[1:]
    return doc
}
func (hlq *HostLoadQueue)Size() int {
    return len(hlq.hosts)
}

type HostLoader struct {
    hostMap map[string]HostLoadQueue
}
func (hl *HostLoader)Travel(f func(*proto.CrawlDoc)) {
    // release hostMap
    // call func f for avaliable doc
}
