package base

import (
	"mustard/crawl/proto"
	"mustard/base/string_util"
	"strings"
)

const (
	kMaxValidUrlLength = 512
)

func GetHostName(doc *proto.CrawlDoc) string {
	if string_util.IsEmpty(doc.CrawlParam.FakeHost) {
		return doc.CrawlParam.FetchHint.Host
	}
	return doc.CrawlParam.FakeHost
}

// TODO call this function in where???
func IsInvalidUrl(_url string) bool {
	/*
		1. start with http or https
		2. url len should less then kMaxValidUrlLength
	*/
	if !(strings.HasPrefix(_url, "http://") || strings.HasPrefix(_url, "https://")) {
		return false
	}
	if len(_url) > kMaxValidUrlLength {
		return false
	}
	return true
}