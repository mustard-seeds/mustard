package scheduler
import (
    "strconv"
    pb "mustard/crawl/proto"
    "mustard/base"
    "mustard/base/conf"
    LOG "mustard/base/log"
    "mustard/utils/url_parser"
    "regexp"
    crawl_base "mustard/crawl/base"
    "mustard/base/file"
)
var CONF = conf.Conf
/*
hostload,multifetcher,fake host,receivers

priority
tag,prime second
randomhostload
drop content
store engine
store db,table
request type
*/
type JobDescription struct {
    isUrgent        bool
    primeTag        string
    secondTag       []string
    randomHostLoad  int
    dropContent     bool
    storeEngine     string
    storeDb         string
    storeTable      string
    requestType     int
}

var NormalJobD = JobDescription{
    isUrgent:false,
    primeTag:"n",
    randomHostLoad:0,
    dropContent:false,
    requestType:1,
}

var UrgentJobD = JobDescription{
    isUrgent:true,
    primeTag:"U",
    randomHostLoad:0,
    dropContent:false,
    requestType:1,
}
type ParamFillerMaster struct {
    fillers ParamFillerGroup
    jd  *JobDescription
}
func (m *ParamFillerMaster)RegisterParamFillerGroup(f ParamFillerGroup) {
    m.fillers = f
}
func (m *ParamFillerMaster)RegisterJobDescription(jd *JobDescription) {
    m.jd = jd
}
func (m *ParamFillerMaster)Init() {
    for _,v := range m.fillers.Fillers() {
        v.Init()
    }
}
func (m *ParamFillerMaster)Fill(doc *pb.CrawlDoc) {
    for _,v := range m.fillers.Fillers() {
        v.Fill(m.jd,doc)
    }
}

type ParamFillerGroup interface {
    Package()
    Fillers() []ParamFiller
}

type DefaultParamFillerGroup struct {
    fillers []ParamFiller
}
func (d *DefaultParamFillerGroup)Fillers() []ParamFiller {
    return d.fillers
}
func (d *DefaultParamFillerGroup)Package() {
    // pay attention the sequence
    // FakeHostParamFiller & HostLoadParamFiller & MultiFetcherParamFiller mush use and ensure the sequence
    d.fillers = append(d.fillers,&PrepareParamFiller{})
    d.fillers = append(d.fillers,&FakeHostParamFiller{})
    d.fillers = append(d.fillers,&HostLoadParamFiller{})
    d.fillers = append(d.fillers,&MultiFetcherParamFiller{})
    d.fillers = append(d.fillers,&ReceiverParamFiller{})
    d.fillers = append(d.fillers,&TagParamFiller{})
}


type ParamFiller interface {
    Init()
    Fill(*JobDescription, *pb.CrawlDoc)
}

// prepareParamFiller should the first one
type PrepareParamFiller struct {
}
func (p *PrepareParamFiller)Init() {
}
func (p *PrepareParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    base.CHECK(doc.RequestUrl != "", "Doc Request url not filled")
    // normalize request_url, fill url,host,path ...
    if doc.GetCrawlParam() == nil {
        doc.CrawlParam = &pb.CrawlParam{}
    }
    if doc.CrawlParam.GetFetchHint() == nil {
        doc.CrawlParam.FetchHint = &pb.FetchHint{}
    }
    // fill url
    doc.Url = url_parser.NormalizeUrl(doc.RequestUrl)
    doc.CrawlParam.FetchHint.Host = url_parser.GetURLObj(doc.Url).Host
    doc.CrawlParam.FetchHint.Path = url_parser.GetURLObj(doc.Url).Path
}

type FakeHostParamFiller struct {
    fakehost map[string]string
    last_load_time  int64
}
func (f *FakeHostParamFiller)loadFakeHostConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.FakeHostConfigFile)
    result,fresh := crawl_base.LoadConfigWithTwoField("FakeHost", fname, ",")
    if fresh {
        for k, v := range result {
            f.fakehost[k] = v
            LOG.VLog(3).Debugf("Load FakeHost %s : %s", k, v)
        }
    }
}
func (f *FakeHostParamFiller)Init() {
    f.fakehost = make(map[string]string)
    f.loadFakeHostConfigFile()
}
func (f *FakeHostParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    f.loadFakeHostConfigFile()
    for k,v := range f.fakehost {
        r,_ := regexp.Compile(k)
        regexRet := r.FindAllString(doc.CrawlParam.FetchHint.Host, -1)
        if len(regexRet) != 0 {
            doc.CrawlParam.FakeHost = v
        }
    }
}

type HostLoadParamFiller struct {
    hostload map[string]int
    last_load_time  int64
}
func (h *HostLoadParamFiller)loadHostloadConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.HostLoadConfigFile)
    result,fresh := crawl_base.LoadConfigWithTwoField("HostLoad", fname, ",")
    if fresh {
        for k,v := range result {
            hl,err := strconv.Atoi(v)
            if err != nil {
                LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k,v)
                return
            }
            h.hostload[k] = hl
            LOG.VLog(3).Debugf("Load HostLoad %s : %d", k, hl)
        }
    }
}
func (h *HostLoadParamFiller)Init() {
    h.hostload = make(map[string]int)
    h.loadHostloadConfigFile()
}
func (h *HostLoadParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    h.loadHostloadConfigFile()  // reload
    host := crawl_base.GetHostName(doc)
    hl := *CONF.Crawler.DefaultHostLoad
    thl,present := h.hostload[host]
    if present {
        hl = thl
    }
    doc.CrawlParam.Hostload = int32(hl)
}

type MultiFetcherParamFiller struct {
    multifetcher map[string]int
    last_load_time  int64
}
func (f *MultiFetcherParamFiller)loadMultiFetcherConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.MultiFetcherConfigFile)
    result,fresh := crawl_base.LoadConfigWithTwoField("HostLoad", fname, ",")
    if fresh {
        for k,v := range result {
            hl,err := strconv.Atoi(v)
            if err != nil {
                LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k,v)
                return
            }
            f.multifetcher[k] = hl
            LOG.VLog(3).Debugf("Load Multifetcher %s : %d", k, hl)
        }
    }
}
func (f *MultiFetcherParamFiller)Init() {
    f.multifetcher = make(map[string]int)
    f.loadMultiFetcherConfigFile()
}
func (f *MultiFetcherParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    f.loadMultiFetcherConfigFile()
    host := crawl_base.GetHostName(doc)
    mf := 1
    thl,present := f.multifetcher[host]
    if present {
        mf = thl
    }
    doc.CrawlParam.FetcherCount = int32(mf)
}

type ReceiverParamFiller struct {
    receivers map[string]*pb.ConnectionInfo
    last_load_time int64
}
func (f *ReceiverParamFiller)loadReceiverConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.ReceiversConfigFile)
    result,fresh := crawl_base.LoadConfigWithTwoField("HostLoad", fname, ":")
    if fresh {
        for k,v := range result {
            hl,err := strconv.Atoi(v)
            if err != nil {
                LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k,v)
                return
            }
            f.receivers[k+":"+v] = &pb.ConnectionInfo{Host:k,Port:int32(hl)}
            LOG.VLog(3).Debugf("Load receivers %s : %d", k, hl)
        }
    }
}
func (f *ReceiverParamFiller)Init() {
    f.receivers = make(map[string]*pb.ConnectionInfo)
    f.loadReceiverConfigFile()
}
func (f *ReceiverParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    f.loadReceiverConfigFile()
    for _,v := range f.receivers {
        doc.CrawlParam.Receivers = append(doc.CrawlParam.Receivers, v)
    }
}

type TagParamFiller struct {
}
func (h *TagParamFiller)Init() {
}
func (h *TagParamFiller)Fill(jd *JobDescription, doc *pb.CrawlDoc) {
    doc.CrawlParam.Pri = pb.Priority_NORMAL
    if jd.isUrgent {
        doc.CrawlParam.Pri = pb.Priority_URGENT
    }
    doc.CrawlParam.PrimaryTag = jd.primeTag
    for _,v := range jd.secondTag {
        doc.CrawlParam.SecondaryTag = append(doc.CrawlParam.SecondaryTag,v)
    }
    doc.CrawlParam.RandomHostload = int32(jd.randomHostLoad)
    doc.CrawlParam.DropContent = jd.dropContent
    doc.CrawlParam.Rtype = pb.RequestType(jd.requestType)
    // storage
    doc.CrawlParam.StoreEngine = jd.storeEngine
    doc.CrawlParam.StoreDb = jd.storeDb
    doc.CrawlParam.StoreTable = jd.storeTable
}