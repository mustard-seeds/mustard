
default: build

build: 
	mkdir -p ./bin
# for example
#	protoc --go_out=plugins=grpc:. example/proto/example.proto
#	go build -o bin/example_service ./example/service
#	go build -o bin/example_client ./example/client
#	go build -o bin/example_client_with_http ./example/client/http_main
# for crawler
#	protoc --go_out=plugins=grpc:. crawl/proto/crawldoc.proto
	go build -o bin/dispatcher_main ./crawl/dispatcher/main
	go build -o bin/fetcher_main ./crawl/fetcher/main
	go build -o bin/file-scheduler_main ./crawl/scheduler/file_scheduler_main
test:
	go test ./base/...
	go test ./crawl/...
	go test ./utils/...
doc:
	godoc -http=:6060 -index
fmt:
	go fmt ./...
clean:
	rm ./bin/*_main
