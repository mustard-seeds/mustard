package main

import (
    pb "mustard/example/proto"
    "mustard/internal/google.golang.org/grpc"
    "mustard/internal/golang.org/x/net/context"
    conf "mustard/base/conf"
    LOG "mustard/base/log"
    "fmt"
    "mustard/base/proto_util"
    "mustard/internal/google.golang.org/grpc/credentials"
)

var CONF = conf.Conf

func general(client pb.GreetingServiceClient) {
    greet, err := client.Greeting(context.Background(), new(pb.GreetRequest))
    if err != nil {
        LOG.Fatal("%v.GetFeatures(_) = _, %v: ", client, err)
    }
    LOG.Info("SayHello return :" + proto_util.FromProtoToString(greet))
}

func greet(client pb.GreetingServiceClient) {
    var request pb.GreetRequest
    request.Hobbies = make([]string, 2)
    request.Hobbies[0] = "swimming"
    request.Hobbies[1] = "swimming"
    request.Keyword = make(map[string]int32)
    request.Keyword["c"] = 1
    request.Person = &pb.Person{"Alice", 21}
    greet, err := client.GreetOnce(context.Background(), &request)
    if err != nil {
        LOG.Fatal("%v.GetFeatures(_) = _, %v: ", client, err)
    }
    LOG.Info(greet.Content)
    LOG.Info("greet return :" + proto_util.FromProtoToString(greet))
}

func main() {
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
    serverAddr := fmt.Sprintf("127.0.0.1:%d", *CONF.Example.Port)
    conn, err := grpc.Dial(serverAddr, opts...)
    if err != nil {
        LOG.Fatalf("fail to dial: %v", err)
    }
    defer conn.Close()
    client := pb.NewGreetingServiceClient(conn)
    general(client)
    greet(client)
}
