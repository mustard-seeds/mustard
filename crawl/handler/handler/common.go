package handler

import (
	"reflect"
)

var CrawlTaskRegistry = make(map[string]CrawlTask)

func registerCrawlTaskType(task CrawlTask) {
	t := reflect.TypeOf(task).Elem()
	CrawlTaskRegistry[t.Name()] = task
}

func GetCrawlHandlerByName(name string) CrawlTask {
	return CrawlTaskRegistry[name]
}
