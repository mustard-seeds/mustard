package main

import (
	"mustard/crawl/dispatcher"
	"fmt"
	"net"
	LOG "mustard/base/log"
	"mustard/internal/google.golang.org/grpc/credentials"
	"mustard/base/conf"
	pb "mustard/crawl/proto"
	"mustard/internal/google.golang.org/grpc"
	"mustard/utils/babysitter"
)
var CONF = conf.Conf

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *CONF.Crawler.DispatcherPort))
	if err != nil {
		LOG.Fatalf("Dispatcher Failed to listen:%v", err)
	} else {
		LOG.Infof("Dispatcher Listen on port: %d", *CONF.Crawler.DispatcherPort)
	}
	var opts []grpc.ServerOption
	if *CONF.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(*CONF.CertFile, *CONF.KeyFile)
		if err != nil {
			LOG.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	disp := new(dispatcher.Dispatcher)
	disp.Init()

	var http_server babysitter.MonitorServer
	http_server.Init()
	http_server.AddMonitor(disp)
	go http_server.Serve(*CONF.Crawler.HttpPort)

	pb.RegisterCrawlServiceServer(grpcServer, disp)
	grpcServer.Serve(lis)
}
