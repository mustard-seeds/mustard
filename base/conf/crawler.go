package conf

import (
    "flag"
)

type CrawlerType struct {
    ChannelBufSize      *int
    CrawlHandlerChain   *string
    CrawlInputHandler   *string
    HostLoaderQueueSize *int
    HostLoaderReleaseRatio *float64
    FetchConnectionNum  *int
    CrawlRequestPort    *int
    CrawlRequestHealthyRatio *float64
}

var CrawlerConf = CrawlerType{
    ChannelBufSize: flag.Int("channel_buffer_size", 100, "channel buffer size"),
    CrawlHandlerChain: flag.String("crawl_handler_chain", "FetchHandler;DocHandler;StorageHandler","handler chain, split by ;"),
    CrawlInputHandler:flag.String("crawl_input_processor", "RequestProcessor","input processors,split by ;"),
    HostLoaderQueueSize:flag.Int("host_load_queue_size", 20, "queue size for each host"),
    HostLoaderReleaseRatio:flag.Float64("host_load_release_ratio",0.6,"release ratio vacancy rate"),
    FetchConnectionNum:flag.Int("fetch_connection_number",10,"url fetch connection number"),
    CrawlRequestPort:flag.Int("crawl_request_port", 9010, "grpc port"),
    CrawlRequestHealthyRatio:flag.Float64("crawl_request_healthy_ratio", 0.9, " healthy raito"),
}
