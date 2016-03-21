package handler

import (
	"fmt"
	LOG "mustard/base/log"
	"mustard/base/string_util"
	"mustard/base/time_util"
	"mustard/crawl/base"
	pb "mustard/crawl/proto"
	"mustard/internal/golang.org/x/net/context"
	"mustard/internal/google.golang.org/grpc"
	"mustard/internal/google.golang.org/grpc/credentials"
	"net"
)

type RequestProcessor struct {
	CrawlHandler
}

func (request *RequestProcessor) Feed(ctx context.Context, docs *pb.CrawlDocs) (*pb.CrawlResponse, error) {
	t1 := time_util.GetTimeInMs()
	healthy := request._isHealthy()
	if healthy {
		for _, doc := range docs.Docs {
			LOG.VLog(4).Debugf("Get doc %s,type:%d", doc.RequestUrl, doc.CrawlParam.Rtype)
			LOG.VLog(4).Debugf("DumpDoc:\n%s", base.DumpCrawlDoc(doc))
			//request.Output(doc)
			request.output_chan <- doc
			request.CrawlHandler.process_num++
			request.accept_num++
		}
	}
	LOG.VLog(3).Debugf("Feed using time %d ms for %d records.", (time_util.GetTimeInMs() - t1), len(docs.Docs))
	return &pb.CrawlResponse{
		Ok:  healthy,
		Ret: int64(*CONF.Crawler.ChannelBufSize - len(request.output_chan)),
	}, nil
}
func (request *RequestProcessor) _isHealthy() bool {
	return float64(len(request.output_chan))/float64(*CONF.Crawler.ChannelBufSize) < *CONF.Crawler.CrawlRequestHealthyRatio
}
func (request *RequestProcessor) IsHealthy(ctx context.Context, r *pb.CrawlRequest) (*pb.CrawlResponse, error) {
	LOG.VLog(4).Debugf("From %s Call IsHealthy", r.Request)
	return &pb.CrawlResponse{
		Ok:  request._isHealthy(),
		Ret: int64(*CONF.Crawler.ChannelBufSize - len(request.output_chan)),
	}, nil
}
func (h *RequestProcessor) Init() bool {
	LOG.VLog(3).Debug("RequestProcessor Init Finish")
	return true
}
func (h *RequestProcessor) Status(s *string) {
	h.CrawlHandler.Status(s)
	string_util.StringAppendF(s, "[(%t)%d/%d-%g]", h._isHealthy(), len(h.output_chan),
		*CONF.Crawler.ChannelBufSize, *CONF.Crawler.CrawlRequestHealthyRatio)
}
func (request *RequestProcessor) Run(p CrawlProcessor) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *CONF.Crawler.CrawlRequestPort))
	if err != nil {
		LOG.Fatalf("Failed to listen:%v", err)
	} else {
		LOG.Infof("listen on port: %d", *CONF.Crawler.CrawlRequestPort)
	}
	var opts []grpc.ServerOption

	if *CONF.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(*CONF.CertFile, *CONF.KeyFile)
		if err != nil {
			LOG.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	LOG.Infof("Start Request Processor at %d", *CONF.Crawler.CrawlRequestPort)
	grpcServer := grpc.NewServer(opts...)
	// grpc server should set this pointer/ self to serve
	pb.RegisterCrawlServiceServer(grpcServer, request)
	grpcServer.Serve(lis)
}

// use for create instance from a string
func init() {
	registerCrawlTaskType(&RequestProcessor{})
}
