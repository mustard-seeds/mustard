package handler

import (
    "net"
    "fmt"
    LOG "mustard/base/log"
    pb "mustard/crawl/proto"
    "mustard/internal/golang.org/x/net/context"
    "mustard/internal/google.golang.org/grpc"
    "mustard/internal/google.golang.org/grpc/credentials"
)

type RequestProcessor struct {
    CrawlHandler
}
func (request *RequestProcessor)Feed(ctx context.Context, docs *pb.CrawlDocs) (*pb.CrawlResponse, error) {
    healthy := request._isHealthy()
    if healthy {
        for _,doc := range docs.Docs {
            request.Output(doc)
        }
    }
    return &pb.CrawlResponse{
        Ok:healthy,
        Ret:int64(*CONF.Crawler.ChannelBufSize - len(request.output_chan)),
    },nil
}
func (request *RequestProcessor)_isHealthy() bool {
    return float64(len(request.output_chan))/float64(*CONF.Crawler.ChannelBufSize) > *CONF.Crawler.CrawlRequestHealthyRatio
}
func (request *RequestProcessor)IsHealthy(ctx context.Context, r *pb.CrawlRequest) (*pb.CrawlResponse, error) {
    return &pb.CrawlResponse{
        Ok:request._isHealthy(),
        Ret:int64(*CONF.Crawler.ChannelBufSize - len(request.output_chan)),
    },nil
}
func (request *RequestProcessor)Run(p CrawlProcessor) {
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
    pb.RegisterCrawlServiceServer(grpcServer, new(RequestProcessor))
    grpcServer.Serve(lis)
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&RequestProcessor{})
}