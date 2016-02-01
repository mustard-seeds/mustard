package dispatcher

import (
    LOG "mustard/base/log"
    "mustard/base"
    "mustard/base/file"
    "strings"
    "mustard/base/conf"
    "mustard/internal/google.golang.org/grpc"
    "mustard/internal/golang.org/x/net/context"
    "mustard/internal/google.golang.org/grpc/credentials"
    pb "mustard/crawl/proto"
    "mustard/base/string_util"
    "mustard/base/hash"
    "sync"
    "mustard/base/time_util"
    "errors"
    "strconv"
)
var CONF = conf.Conf

const (
    kMaxBatchFeedSize int = 100
    kInvalidCrawlerID uint32 = 65535
    kMagicNumber      uint32 = 113
)
// send crawldoc to target server
// dispatch as:  host/domain/url/random
// remotes:  crawldoc.receivers or configfile.
type CrawlerFeeder struct {
    host string
    port int
    connected bool  // default false
    crawldocs pb.CrawlDocs  // send to fetcher
    docCache []*pb.CrawlDoc  // cache for each fetcher
    client pb.CrawlServiceClient
    // statistic
    process_urls int
    queuefull_urls int
}
func (cf *CrawlerFeeder)PendingUrls() int {
    return len(cf.crawldocs.Docs) + len(cf.docCache)
}
func (cf *CrawlerFeeder)ProcessedUrls() int {
    return cf.process_urls
}
func (cf *CrawlerFeeder)QueuefullUrls() int {
    return cf.queuefull_urls
}

func (cf *CrawlerFeeder)ShouldOutput(doc *pb.CrawlDoc) bool {
    if (cf.PendingUrls() < *CONF.Crawler.FeederMaxPending) {
        return true
    }
    if (doc.CrawlParam.Pri == pb.Priority_URGENT) {
        return true
    } else if (doc.CrawlParam.Pri == pb.Priority_NORMAL) {
        return false
    } else {
        LOG.Warningf("Unknow Priority:%s, pri:%d",doc.Url,doc.CrawlParam.Pri)
        return false
    }
}

func (cf *CrawlerFeeder)AddFeed(doc *pb.CrawlDoc) bool{
    if (cf.ShouldOutput(doc)) {
        cf.docCache = append(cf.docCache, doc)
        return true
    } else {
        cf.queuefull_urls++
        return false
    }
}
func (cf *CrawlerFeeder)Flush() {
    for (len(cf.crawldocs.Docs) < kMaxBatchFeedSize && len(cf.docCache) > 0) {
        doc := cf.docCache[0]
        doc.CrawlRecord.Fetcher.Host = cf.host
        doc.CrawlRecord.Fetcher.Port = int32(cf.port)
        cf.crawldocs.Docs = append(cf.crawldocs.Docs, doc)
        cf.docCache = cf.docCache[1:]
    }
    if (len(cf.crawldocs.Docs) == 0) {
        return
    }
    response,err := cf.client.Feed(context.Background(), &cf.crawldocs)
    if (err != nil) {
        LOG.Errorf("Fail to Feed for Exception %s:%d",cf.host,cf.port)
        cf.connected = false
        return
    }
    if (response.Ok == false) {
        LOG.Errorf("Fail to Feed for UnHealthy %s:%d", cf.host,cf.port)
        cf.connected = false
        return
    }
    cf.process_urls += len(cf.crawldocs.Docs)
    cf.crawldocs.Reset()
}
func (cf *CrawlerFeeder)Connect() bool {
    if (!cf.connected) {
        var opts []grpc.DialOption
        if *CONF.UseTLS {
            var sn string
            if *CONF.ServerHostOverride != "" {
                sn = *CONF.ServerHostOverride
            }
            var creds credentials.TransportAuthenticator
            if *CONF.CaFile != "" {
                var err error
                creds, err = credentials.NewClientTLSFromFile(*CONF.CaFile, sn)
                if err != nil {
                    LOG.Fatalf("Failed to create TLS credentials %v", err)
                }
            } else {
                creds = credentials.NewClientTLSFromCert(nil, sn)
            }
            opts = append(opts, grpc.WithTransportCredentials(creds))
        } else {
            opts = append(opts, grpc.WithInsecure())
        }
        var serverAddr string
        string_util.StringAppendF(&serverAddr,"%s:%d",*CONF.Crawler.DispatcherHost,*CONF.Crawler.DispatcherPort)
        conn,err := grpc.Dial(serverAddr, opts...)
        if err != nil {
            LOG.Errorf("fail to dial: %v", err)
            conn.Close()
        }
        cf.client = pb.NewCrawlServiceClient(conn)
        cf.connected = true
    }
    return cf.connected
}

func (cf *CrawlerFeeder)IsHealthy() bool{
    if (cf.Connect()) {
        response,err := cf.client.IsHealthy(context.Background(),&pb.CrawlRequest{Request:"Dispatch"})
        if (err != nil) {
            cf.connected = false
            return false
        }
        return response.Ok
    }
    return false
}

func (cf *CrawlerFeeder)IsConnected() bool {
    return cf.Connect()
}


type CrawlerFeederGroup struct {
    liveFeeders map[uint32]bool
    deadFeeders map[uint32]bool
    feeders map[uint32]*CrawlerFeeder
}

type Dispatcher struct {
    feeders *CrawlerFeederGroup
    sync.RWMutex
}
func (d *Dispatcher)SelectCrawler(doc *pb.CrawlDoc) uint32{
    if (len(d.feeders.liveFeeders) == 0) {
        return kInvalidCrawlerID;
    }
    url := doc.RequestUrl
    // TODO: check CrawlParam nil or not?
    host := doc.CrawlParam.FetchHint.Host
    var key uint32
    if (*CONF.Crawler.DispatchAs == "host") {
        if (doc.CrawlParam.FakeHost != "") {
            host = doc.CrawlParam.FakeHost
        }
        key = hash.FingerPrint32(host)
    } else {
        key = hash.FingerPrint32(url)
    }
    var base_index uint32 = (key * kMagicNumber) % uint32(len(d.feeders.feeders))
    var offset uint32 = hash.FingerPrint32(url) % uint32(doc.CrawlParam.FetcherCount)  // multi fetcher
    var crawlerId uint32 = (base_index + offset) % uint32(len(d.feeders.feeders))
    var live_crawler_id uint32 = crawlerId

    for (true) {
        _,present := d.feeders.deadFeeders[live_crawler_id]
        if present {
            break
        }
        live_crawler_id++
        live_crawler_id = live_crawler_id % uint32(len(d.feeders.feeders))
        if (live_crawler_id == crawlerId) {
            return kInvalidCrawlerID
        }
    }
    return live_crawler_id
}
func (d *Dispatcher)Flush() {
    t1 := time_util.GetTimeInMs()
    d.Lock()
    defer d.Unlock()
    for k,_ := range d.feeders.liveFeeders {
        if (d.feeders.feeders[k].IsHealthy()) {
            d.feeders.feeders[k].Flush()
        }
    }
    LOG.VLog(3).Debugf("Flush using time %d ms.",time_util.GetTimeInMs() - t1)
}
func (d *Dispatcher)UpdateCrawlerStatus() {
    t1 := time_util.GetTimeInMs()
    d.Lock()
    defer d.Unlock()
    liveF := make(map[uint32]bool)
    deadF := make(map[uint32]bool)
    for k,v := range d.feeders.feeders {
        if (v.IsConnected()) {
            liveF[k] = true
        } else {
            deadF[k] = true
        }
    }
    d.feeders.liveFeeders = liveF
    d.feeders.deadFeeders = deadF
    LOG.VLog(3).Debugf("UpdateCrawlerStatus using time %d ms.",time_util.GetTimeInMs() - t1)
}
func (d *Dispatcher)Feed(ctx context.Context, docs *pb.CrawlDocs) (*pb.CrawlResponse,error) {
    t1 := time_util.GetTimeInMs()
    d.Lock()
    defer d.Unlock()
    for _,v := range docs.Docs {
        var crawlerid uint32 = d.SelectCrawler(v)
        LOG.VLog(3).Debugf("New Url:%s, RequestType:%d,CrawlerId:%d",v.RequestUrl,crawlerid)
        if (crawlerid != kInvalidCrawlerID) {
            if (!d.feeders.feeders[crawlerid].AddFeed(v)) {
                // TODO: SITEQUEUEFULL send to fetcher's handler directly
                LOG.Errorf("NOT IMPLEMENT, send to fetcher %s",v.RequestUrl)
            }
        }
    }
    LOG.VLog(3).Debugf("Feed using time %d ms for %d records.",(time_util.GetTimeInMs() - t1),len(docs.Docs))
    return &pb.CrawlResponse{Ok:true},nil
}
func (d *Dispatcher)IsHealthy(ctx context.Context, request *pb.CrawlRequest) (*pb.CrawlResponse, error) {
    d.Lock()
    defer d.Unlock()
    if (float64(len(d.feeders.liveFeeders)) < (float64(len(d.feeders.liveFeeders) + len(d.feeders.deadFeeders))* *CONF.Crawler.DispatchLiveFeederRatio)) {
        LOG.VLog(3).Debugf("Live feeders only %d Not Healthy",len(d.feeders.liveFeeders))
        return &pb.CrawlResponse{Ok:false},errors.New("Live Feeders Too little")
    }
    var pending_urls int = 0
    for _,v := range d.feeders.feeders {
        pending_urls += v.PendingUrls()
    }
    if (pending_urls > *CONF.Crawler.GroupFeederMaxPending) {
        LOG.VLog(3).Debugf("Pending url too much : %d", pending_urls)
        return &pb.CrawlResponse{Ok:false},errors.New("Pending Url Too Much")
    }
    return &pb.CrawlResponse{Ok:true},nil
}
func (d *Dispatcher)Init() {
    d.feeders = &CrawlerFeederGroup{
        liveFeeders:make(map[uint32]bool),
        deadFeeders:make(map[uint32]bool),
        feeders:make(map[uint32]*CrawlerFeeder),
    }
    d.LoadCrawlersFromFile(*CONF.Crawler.CrawlersConfigFile)
    go d.CrawlFeederLoop()
    // start rpc service at dispatcher_main
}
func (d *Dispatcher)CrawlFeederLoop() {
    for (true) {
        d.UpdateCrawlerStatus()
        d.Flush()
        time_util.Sleep(*CONF.Crawler.DispatchFlushInterval)
    }
}

func (d *Dispatcher)LoadCrawlersFromFile(name string) {
    base.CHECK(file.Exist(name))
    content,_ := file.ReadFileToString(name)
    lines := strings.Split(content, "\n")
    for _,l := range lines {
        if l == "" || strings.HasPrefix(l, "#") {
            continue
        }
        addr :=strings.Split(l,":")
        base.CHECK(len(addr) == 2)
        addrPort,err := strconv.Atoi(addr[1])
        base.CHECK(err == nil)
        d.feeders.feeders[uint32(len(d.feeders.feeders))] = &CrawlerFeeder{
            host:addr[0],
            port:addrPort,
        }
    }
}
