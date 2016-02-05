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
    "mustard/utils/babysitter"
    "time"
)

var CONF = conf.Conf

const (
    kMaxBatchFeedSize int = 100
    kInvalidCrawlerID uint32 = 65535
    kMagicNumber      uint32 = 113
    kFeedSpeedInterval int64 = 60
    kStatusPageColSize int = 8

    kPortStep = 50
)
// dispatch as:  host/domain/url/random
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
    LOG.VLog(4).Debugf("Feeder Try Flush CacheDoc %d to %s:%d", len(cf.docCache), cf.host, cf.port)
    for (len(cf.crawldocs.Docs) < kMaxBatchFeedSize && len(cf.docCache) > 0) {
        doc := cf.docCache[0]
        if doc.CrawlRecord.GetFetcher() == nil {
            doc.CrawlRecord.Fetcher = &pb.ConnectionInfo{}
        }
        doc.CrawlRecord.Fetcher.Host = cf.host
        doc.CrawlRecord.Fetcher.Port = int32(cf.port)
        cf.crawldocs.Docs = append(cf.crawldocs.Docs, doc)
        cf.docCache = cf.docCache[1:]
    }
    if (len(cf.crawldocs.Docs) == 0) {
        return
    }
    LOG.VLog(4).Debugf("Feeder Prepare Flush %d docs to %s:%d", len(cf.crawldocs.Docs), cf.host, cf.port)
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
        opts = append(opts,grpc.WithTimeout(time.Second*time.Duration(*CONF.Crawler.ConnectionTimeout)))
        opts = append(opts,grpc.WithBlock()) // grpc should with block...
        var serverAddr string
        string_util.StringAppendF(&serverAddr,"%s:%d",cf.host,cf.port)
        conn,err := grpc.Dial(serverAddr, opts...)
        if err != nil {
            LOG.Errorf("fail to dial %s: %v",serverAddr, err)
            cf.connected = false
        } else {
            cf.client = pb.NewCrawlServiceClient(conn)
            cf.connected = true
            LOG.Infof("Connect Feeder %s",serverAddr)
        }
    }
    return cf.connected
}

func (cf *CrawlerFeeder)IsHealthy() bool{
    if (cf.Connect()) {
        // TODO, if dispatcher call itself, will deadlock
        response,err := cf.client.IsHealthy(context.Background(),&pb.CrawlRequest{Request:"Dispatch"})
        if (err != nil) {
            LOG.VLog(2).Debugf("Connect %s:%d Error %s",cf.host,cf.port,err.Error())
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
    // statistic
    last_counter_time int64
    input_speed int64
    input_counter int64
    start_time_str string
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
        if !present {
            break
        }
        // if live_crawler_id is in deadFeeders, it should shift next
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
        LOG.VLog(3).Debugf("New Url:%s, RequestType:%d,CrawlerId:%d",v.RequestUrl,v.CrawlParam.Rtype,crawlerid)
        if (crawlerid != kInvalidCrawlerID) {
            if (!d.feeders.feeders[crawlerid].AddFeed(v)) {
                // TODO: SITEQUEUEFULL send to fetcher's handler directly
                LOG.Errorf("NOT IMPLEMENT, send to fetcher %s",v.RequestUrl)
            }
        }
    }
    d.input_counter += int64(len(docs.Docs))
    LOG.VLog(3).Debugf("Feed using time %d ms for %d records.",(time_util.GetTimeInMs() - t1),len(docs.Docs))
    return &pb.CrawlResponse{Ok:true},nil
}
func (d *Dispatcher)IsHealthy(ctx context.Context, request *pb.CrawlRequest) (*pb.CrawlResponse, error) {
    LOG.VLog(4).Debugf("From %s Call Dispatcher IsHealthy", request.Request)
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
    fname := file.GetConfFile(*CONF.Crawler.CrawlersConfigFile)
    d.LoadCrawlersFromFile(fname)
    go d.CrawlFeederLoop()
    // start rpc service at dispatcher_main
    d.start_time_str = time_util.GetReadableTimeNow()
}
func (d *Dispatcher)CrawlFeederLoop() {
    for (true) {
        d.UpdateCrawlerStatus()
        d.Flush()
        time_util.Sleep(*CONF.Crawler.DispatchFlushInterval)
    }
}

func (d *Dispatcher)LoadCrawlersFromFile(name string) {
    file.FileLineReader(name,"#",func(line string){
        addr :=strings.Split(line,":")
        base.CHECK(len(addr) == 2, "Parse Addr Fail. %s", line)
        addrPort,err := strconv.Atoi(addr[1])
        base.CHECK(addrPort != *CONF.Crawler.DispatcherPort, "Dispatch itself, %s:%d",addr[0],addrPort)
        base.CHECK(err == nil,"atoi error %s", err)
        d.feeders.feeders[uint32(len(d.feeders.feeders))] = &CrawlerFeeder{
            host:addr[0],
            port:addrPort,
            connected:false,
        }
    })
}
func (d *Dispatcher)MonitorReport(result *babysitter.MonitorResult) {
    var info string
    var process_urls,pending_urls,queuefull_urls int
    var interval = time_util.GetCurrentTimeStamp() - d.last_counter_time
    if (interval >= kFeedSpeedInterval) {
        d.input_speed = d.input_counter / interval
        LOG.VLog(3).Debugf("InputCounter:%d,interval:%d,Speed:%d",d.input_counter,interval,d.input_speed)
        d.input_counter = 0
        d.last_counter_time  = time_util.GetCurrentTimeStamp()
    }
    string_util.StringAppendF(&info, "InputSpeed: %d/sec<br>", d.input_speed)
    td_list := []string{}
    for k,_ := range d.feeders.liveFeeders {
        var tds string
        d.fillTDString(d.feeders.feeders[k], true, &tds)
        process_urls += d.feeders.feeders[k].process_urls
        pending_urls += d.feeders.feeders[k].PendingUrls()
        queuefull_urls += d.feeders.feeders[k].queuefull_urls
        td_list = append(td_list, tds)
    }
    for k,_ := range d.feeders.deadFeeders {
        var tds string
        d.fillTDString(d.feeders.feeders[k], false, &tds)
        process_urls += d.feeders.feeders[k].process_urls
        pending_urls += d.feeders.feeders[k].PendingUrls()
        queuefull_urls += d.feeders.feeders[k].queuefull_urls
        td_list = append(td_list, tds)
    }


    string_util.StringAppendF(&info,
        "<style type=\"text/css\">"+
        " table {font-size: 80%%;}"+
        "</style>"+
        "<b> Crawler Dispatcher Summary (start at %s) </b><pre>"+
        "<key>Live feeders  </key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>"+
        "<key>Dead feeders  </key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>"+
        "<key>Processed Urls</key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>"+
        "<key>Pending Urls  </key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>"+
        "<key>QueueFull Urls</key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>"+
        "</pre>",
        d.start_time_str,
        len(d.feeders.liveFeeders),
        len(d.feeders.deadFeeders),
        process_urls,
        pending_urls,
        queuefull_urls)
    info += "<table border=1>"
    for i:=0;i < len(td_list);i++ {
        info += "<tr>"
        j := 0
        for ;i< kStatusPageColSize;j++ {
            if (i + j == len(td_list)) {
                break
            }
            info += td_list[i+j]
        }
        info += "</tr>"
        i += j - 1
    }
    info += "</table>";
    result.AddString(info)
}
func (d *Dispatcher)MonitorReportHealthy() error {
    _,_e := d.IsHealthy(context.Background(),&pb.CrawlRequest{Request:"ItselfCheck"})
    return _e
}
func (d *Dispatcher)fillTDString(feeder *CrawlerFeeder, alive bool, tds *string) {
    var status string
    if alive {
        status = "<font color=green><em>live</em></font>"
    } else {
        status = "<font color=red><em>dead</em></font>";
    }
    string_util.StringAppendF(tds,
        "<td><div><a href=http://%s:%d/statusi>%s:%d</a>%s<pre>" +
        "<key>Processed Url </key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>" +
        "<key>Pending   Urls</key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>" +
        "<key>QueueFull Urls</key>&nbsp;&nbsp;:&nbsp;&nbsp;<value>%d</value><br>" +
        "</pre></div></td>",
        feeder.host,
        feeder.port + kPortStep,
        feeder.host,
        feeder.port,
        status,
        feeder.process_urls,
        feeder.PendingUrls(),
        feeder.queuefull_urls)
}
