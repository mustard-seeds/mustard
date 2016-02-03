
default: build

build: 
	mkdir -p ./bin
	go build -o bin/dispatcher_main ./crawl/dispatcher/main
	go build -o bin/fetcher_main ./crawl/fetcher/main
	go build -o bin/file_scheduler_main ./crawl/scheduler/file_scheduler_main
test:
	go test ./base/...
	go test ./crawl/...
	go test ./utils/...
clean:
	rm ./bin/*_main
