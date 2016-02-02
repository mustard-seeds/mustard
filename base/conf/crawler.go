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
    CrawlersConfigFile  *string
    DefaultHostLoad     *int
    HostLoadConfigFile  *string
    ReceiversConfigFile  *string
    MultiFetcherConfigFile  *string
    FakeHostConfigFile  *string
    FeederMaxPending    *int
    GroupFeederMaxPending    *int
    DispatcherHost      *string
    DispatcherPort      *int
    DispatchAs          *string
    DispatchLiveFeederRatio *float64
    DispatchFlushInterval  *int
    HttpPort            *int
    ConnectionTimeout   *int
    ConfigFileReloadInterval *int
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
    CrawlersConfigFile:flag.String("crawlers_config_file","etc/crawl/crawlers.config", "fetcher config file, ip:port each line"),
    HostLoadConfigFile:flag.String("hostload_config_file", "etc/crawl/hostload.config", "hostload config file"),
    ReceiversConfigFile:flag.String("receivers_config_file", "etc/crawl/receivers.config", "receivers config file"),
    MultiFetcherConfigFile:flag.String("multifetcher_config_file", "etc/crawl/multifetcher.config", "multi fetcher config file"),
    FakeHostConfigFile:flag.String("fakehost_config_file", "etc/crawl/fakehost.config", "multi fetcher config file"),
    FeederMaxPending:flag.Int("feeder_max_pendings",100,"feeder max pending for dispatcher"),
    GroupFeederMaxPending:flag.Int("group_feeder_max_pendings",5000,"feeder max pending for dispatcher"),
    DispatcherHost:flag.String("dispatcher_host","127.0.0.1","dispatcher host"),
    DispatcherPort:flag.Int("dispatcher_port",9000,"dispatcher port"),
    DispatchAs:flag.String("dispatch_as","host","host or url, dispatch as"),
    DispatchLiveFeederRatio:flag.Float64("live_feeder_ratio", 0, "dispatcher live feeder ratio"),
    DispatchFlushInterval:flag.Int("dispatch_flush_interval", 10,"dispatch flush interval"),
    HttpPort:flag.Int("http_port",9900,"http port"),
    ConnectionTimeout:flag.Int("connection_timeout",2,"connection timeout"),
    DefaultHostLoad:flag.Int("default_hostload",5,"default host load"),
    ConfigFileReloadInterval:flag.Int("config_file_reload_interval",1800, "config file reload interval"),
}
