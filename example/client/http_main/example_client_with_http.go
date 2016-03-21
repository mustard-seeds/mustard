package main

import (
	"fmt"
	"mustard/base/conf"
	LOG "mustard/base/log"
	"mustard/base/proto_util"
	pb "mustard/example/proto"
	"mustard/internal/github.com/gorilla/mux"
	"mustard/internal/golang.org/x/net/context"
	"mustard/internal/google.golang.org/grpc"
	"mustard/internal/google.golang.org/grpc/credentials"
	"net/http"
)

var CONF = conf.Conf

var client pb.GreetingServiceClient

func Greetings(w http.ResponseWriter, r *http.Request) {
	greet, err := client.Greeting(context.Background(), new(pb.GreetRequest))
	if err != nil {
		LOG.Fatal("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	LOG.Info("SayHello return :" + proto_util.FromProtoToString(greet))
	w.Write([]byte(proto_util.FromProtoToString(greet)))
}

func Greeting(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(proto_util.FromProtoToString(greet)))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Greetings)
	r.HandleFunc("/{name}", Greeting)

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
	client = pb.NewGreetingServiceClient(conn)

	LOG.Info("Starting REST server")
	http.ListenAndServe(":8080", r)
}
