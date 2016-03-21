package scheduler

import (
	LOG "mustard/base/log"
	"mustard/base/string_util"
	"mustard/base/time_util"
	pb "mustard/crawl/proto"
	"mustard/internal/golang.org/x/net/context"
	"mustard/internal/google.golang.org/grpc"
	"mustard/internal/google.golang.org/grpc/credentials"
	"time"
)

const (
	kCrawlDocSendRetryTimes = 5
	kCrawlDocSendRetryInterval = 1
)

// send to CrawlServiceServer
type CrawlDocSender struct {
	host      string
	port      int
	Connected bool
	client    pb.CrawlServiceClient

	// statistic
	timestamp            int64
	url_sent_in_sec      int
	max_url_send_per_sec int
	send_speed           int

	LastReconnectTimeStamp int64
}

func (s *CrawlDocSender) SetMaxSpeed(speed int) {
	s.max_url_send_per_sec = speed
}
func (s *CrawlDocSender) Init() {
	s.Connect()
}
func (s *CrawlDocSender) Flush(doc *pb.CrawlDoc) {
	for i := 0; i < kCrawlDocSendRetryTimes;i++ {
		if s.Connect() {
			break
		}
		LOG.VLog(2).Debugf("Connect CrawlServiceServer %s:%d Fail", s.host, s.port)
		time_util.Sleep(kCrawlDocSendRetryInterval)
	}
	if s.Connected == false {
		LOG.Errorf("[FLUSH]Can not send to %s:%d Drop Doc %s", s.host, s.port, doc.String())
		return
	}
	for i := 0; i < kCrawlDocSendRetryTimes;i++ {
		// check healthy
		_, err := s.client.IsHealthy(context.Background(), &pb.CrawlRequest{Request: "CrawlDocSender"})
		if err == nil {
			break
		}
		LOG.VLog(2).Debugf("UnHealthy CrawlServiceServer %s:%d. Reason:%s", s.host, s.port, err.Error())
		s.Connected = false
		time_util.Sleep(kCrawlDocSendRetryInterval * 2)
	}
	if s.Connected == false {
		LOG.Errorf("[FLUSH]NotHealthy %s:%d Drop Doc %s", s.host, s.port, doc.String())
		return
	}

	// Make Sure s.Connected is True before below code.
	now := time_util.GetCurrentTimeStamp()
	if now > s.timestamp {
		s.timestamp = now
		LOG.VLog(1).Debugf("Send url in second %d", s.url_sent_in_sec)
		s.send_speed = s.url_sent_in_sec
	} else if s.url_sent_in_sec >= s.max_url_send_per_sec {
		for now <= s.timestamp {
			time_util.Usleep(10 * 1000)
			now = time_util.GetCurrentTimeStamp()
		}
		LOG.VLog(1).Debugf("Send url Speed %d/Second", s.url_sent_in_sec)
		s.send_speed = s.url_sent_in_sec
		s.url_sent_in_sec = 0
		s.timestamp = now
	}
	doc.CrawlRecord.RequestTime = now
	docs := pb.CrawlDocs{}
	docs.Docs = append(docs.Docs, doc)
	LOG.VLog(2).Debugf("Send %s, host:%s", doc.RequestUrl, doc.CrawlParam.FetchHint.Host)
	_, e := s.client.Feed(context.Background(), &docs)
	if e != nil {
		s.Connected = false
	}
	s.url_sent_in_sec++
}
func (s *CrawlDocSender) Connect() bool {
	if !s.Connected {
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
		opts = append(opts, grpc.WithTimeout(time.Second*time.Duration(*CONF.Crawler.ConnectionTimeout)))
		opts = append(opts, grpc.WithBlock()) // grpc should with block...
		var serverAddr string
		string_util.StringAppendF(&serverAddr, "%s:%d", s.host, s.port)
		conn, err := grpc.Dial(serverAddr, opts...)
		s.LastReconnectTimeStamp = time_util.GetCurrentTimeStamp()
		if err != nil {
			LOG.Errorf("fail to dial %s: %v", serverAddr, err)
			s.Connected = false
		} else {
			LOG.Infof("Connect Feeder %s", serverAddr)
			s.client = pb.NewCrawlServiceClient(conn)
			s.Connected = true
		}
	}
	return s.Connected
}

func NewCrawlDocSender(host string, port int, speed int) *CrawlDocSender {
	return &CrawlDocSender{
		host:                 host,
		port:                 port,
		Connected:            false,
		max_url_send_per_sec: speed,
	}
}
