package handler

import (
    "reflect"
)

var CrawlHandlerRegistry = make(map[string]reflect.Type)

func registerType(typedNil interface{}) {
    t := reflect.TypeOf(typedNil).Elem()
    CrawlHandlerRegistry[t.Name()] = t
}

func GetCrawlHandlerByName(name string) interface{} {
    return reflect.New(CrawlHandlerRegistry[name]).Elem().Interface()
}
