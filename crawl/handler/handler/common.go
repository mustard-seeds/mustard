package handler

import (
    "reflect"
)

var CrawlHandlerRegistry = make(map[string]CrawlHandler)
var CrawlProcessorRegistry = make(map[string]CrawlProcessor)

func registerHandlerType(handler CrawlHandler) {
    t := reflect.TypeOf(handler).Elem()
    CrawlHandlerRegistry[t.Name()] = handler
}
func registerProcessorType(processor CrawlProcessor) {
    t := reflect.TypeOf(processor).Elem()
    CrawlProcessorRegistry[t.Name()] = processor
}

func GetCrawlHandlerByName(name string) CrawlHandler {
    return CrawlHandlerRegistry[name]
}
func GetCrawlProcessorByName(name string) CrawlProcessor {
    return CrawlProcessorRegistry[name]
}
