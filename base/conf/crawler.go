package conf

import (
    "flag"
)

type CrawlerType struct {
    PeriodicInterval    *int
    ChannelBufSize      *int
    CrawlHandlerChain   *string
    CrawlInputHandler   *string
    HostLoaderQueueSize *int
    HostLoaderReleaseRatio *float64
    FetchConnectionNum  *int
}

var CrawlerConf = CrawlerType{
    PeriodicInterval: flag.Int("periodic_interval", 10, "second"),
    ChannelBufSize: flag.Int("channel_buffer_size", 100, "channel buffer size"),
    CrawlHandlerChain: flag.String("crawl_handler_chain", "FetchHandler;DocHandler;StorageHandler","handler chain, split by ;"),
    CrawlInputHandler:flag.String("crawl_input_processor", "RequestProcessor","input processors,split by ;"),
    HostLoaderQueueSize:flag.Int("host_load_queue_size", 20, "queue size for each host"),
    HostLoaderReleaseRatio:flag.Float64("host_load_release_ratio",0.6,"release ratio vacancy rate"),
    FetchConnectionNum:flag.Int("fetch_connection_number",10,"url fetch connection number"),
}
