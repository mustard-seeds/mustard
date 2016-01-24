package main

import (
    "fmt"
    "io"
    "net"
    "time"
    LOG "mustard/base/log"
    "mustard/internal/google.golang.org/grpc/credentials"
    "mustard/base/conf"
    pb "mustard/example/proto"
    "mustard/base/proto_util"
    "mustard/internal/golang.org/x/net/context"
    "mustard/internal/google.golang.org/grpc"
)

var CONF = conf.Conf

// implement of interface GreetingServiceServer
type exampleService struct {
}

func (s *exampleService) Greeting(ctx context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
    time.Sleep(time.Minute)
    LOG.Info("Client request for Greeting")
    return &pb.GreetResponse{"Hello content", "attachment"}, nil
}

func (s *exampleService) GreetOnce(ctx context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
    LOG.Info("get Request GreetOnce :" + proto_util.FromProtoToString(request))
    return &pb.GreetResponse{proto_util.FromProtoToString(request),"attatchment"},nil
}

func (s *exampleService) GreetMulti(stream pb.GreetingService_GreetMultiServer) error {
    for {
        request, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        content := fmt.Sprintf("Request %s", proto_util.FromProtoToString(request))
        if err := stream.Send(&pb.GreetResponse{content, "and your group"}); err != nil {
            return err
        }
    }
}

func main() {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *CONF.Example.Port))
    if err != nil {
        LOG.Fatalf("Failed to listen:%v", err)
    } else {
        LOG.Infof("listen on port: %d", *CONF.Example.Port)
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
    pb.RegisterGreetingServiceServer(grpcServer, new(exampleService))
    grpcServer.Serve(lis)
}