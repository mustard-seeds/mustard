package conf

import (
    "flag"
)

type CrawlerType struct {
    PeriodicInterval    *int
    ChannelBufSize      *int
    CrawlHandlerChain   *string
    CrawlInputHandler   *string
}

var CrawlerConf = CrawlerType{
    PeriodicInterval: flag.Int("periodic_interval", 10, "second"),
    ChannelBufSize: flag.Int("channel_buffer_size", 100, "channel buffer size"),
    CrawlHandlerChain: flag.String("crawl_handler_chain", "DocHandler;StorageHandler","handler chain, split by ;"),
    CrawlInputHandler:flag.String("crawl_input_processor", "RequestProcessor","input processors,split by ;"),
}
