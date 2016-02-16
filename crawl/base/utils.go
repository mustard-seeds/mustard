package base

import (
	"mustard/crawl/proto"
	"mustard/base/string_util"
	"strings"
	"mustard/base/proto_util"
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
func DumpCrawlDoc(doc *proto.CrawlDoc) string {
	docContent := doc.Content
	doc.Content = "..."
	dumpString := proto_util.FromProtoToString(doc)
	doc.Content = docContent
	return dumpString
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
func IsCrawlSuccess(t proto.ReturnType) bool {
	return t == proto.ReturnType_STATUS200 || t == proto.ReturnType_STATUS201
}
func IsPermanentRedirect(t proto.ReturnType) bool {
	return t == proto.ReturnType_STATUS301
}
func IsTemporaryRedirect(t proto.ReturnType) bool {
	return t == proto.ReturnType_STATUS300 ||
			t == proto.ReturnType_STATUS302 ||
			t == proto.ReturnType_STATUS305 ||
			t == proto.ReturnType_STATUS307
}
