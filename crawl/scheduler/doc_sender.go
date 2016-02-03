package scheduler

import (
	"time"
	"mustard/base/string_util"
	pb "mustard/crawl/proto"
	LOG "mustard/base/log"
	"mustard/base/time_util"
	"mustard/internal/golang.org/x/net/context"
	"mustard/internal/google.golang.org/grpc"
	"mustard/internal/google.golang.org/grpc/credentials"

)

// send to CrawlServiceServer
type CrawlDocSender struct {
	host string
	port int
	connected bool
	client pb.CrawlServiceClient
	// statistic
	timestamp int64
	url_sent_in_sec int
	max_url_send_per_sec int
	send_speed int
}
func (s *CrawlDocSender)SetMaxSpeed(speed int) {
	s.max_url_send_per_sec = speed
}
func (s *CrawlDocSender)Init() {
	s.Connect()
}
func (s *CrawlDocSender)Flush(doc *pb.CrawlDoc) {
	for (!s.Connect()) {
		LOG.VLog(2).Debugf("Connect CrawlServiceServer %s:%d Fail", s.host,s.port)
		time_util.Sleep(1)
	}
	for true {
		// check healthy
		_,err := s.client.IsHealthy(context.Background(),&pb.CrawlRequest{Request:"CrawlDocSender"})
		if err == nil {
			break
		}
		LOG.VLog(2).Debugf("UnHealthy CrawlServiceServer %s:%d. Reason:%s", s.host,s.port, err.Error())
		s.connected = false
		time_util.Sleep(2)
	}

	now := time_util.GetCurrentTimeStamp()
	if now > s.timestamp {
		s.timestamp = now
		LOG.VLog(1).Debugf("Send url in second %d",s.url_sent_in_sec)
		s.send_speed = s.url_sent_in_sec
	} else if s.url_sent_in_sec >= s.max_url_send_per_sec {
		for now <= s.timestamp {
			time_util.Usleep(10*1000)
			now = time_util.GetCurrentTimeStamp()
		}
		LOG.VLog(1).Debugf("Send url in Second:%d",s.url_sent_in_sec)
		s.send_speed = s.url_sent_in_sec
		s.url_sent_in_sec = 0
		s.timestamp = now
	}
	doc.CrawlRecord.RequestTime = now
	docs := pb.CrawlDocs{}
	docs.Docs = append(docs.Docs, doc)
	LOG.VLog(2).Debugf("Send %s, host:%s",doc.RequestUrl, doc.CrawlParam.FetchHint.Host)
	_,e := s.client.Feed(context.Background(),&docs)
	if e != nil {
		s.connected = false
	}
	s.url_sent_in_sec++
}
func (s *CrawlDocSender)Connect() bool {
	if (!s.connected) {
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
		opts = append(opts, grpc.WithTimeout(time.Second * time.Duration(*CONF.Crawler.ConnectionTimeout)))
		opts = append(opts, grpc.WithBlock()) // grpc should with block...
		var serverAddr string
		string_util.StringAppendF(&serverAddr, "%s:%d", s.host, s.port)
		conn, err := grpc.Dial(serverAddr, opts...)
		if err != nil {
			LOG.Errorf("fail to dial %s: %v", serverAddr, err)
			s.connected = false
		} else {
			LOG.Infof("Connect Feeder %s", serverAddr)
			s.client = pb.NewCrawlServiceClient(conn)
			s.connected = true
		}
	}
	return s.connected
}

func NewCrawlDocSender(host string, port int, speed int) *CrawlDocSender {
	return &CrawlDocSender{
		host:host,
		port:port,
		connected:false,
		max_url_send_per_sec:speed,
	}
}
