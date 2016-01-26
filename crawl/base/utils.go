package base

import (
	"mustard/crawl/proto"
	"mustard/base/string_util"
)

func GetHostName(doc *proto.CrawlDoc) string {
	if string_util.IsEmpty(doc.CrawlParam.FakeHost) {
		return doc.CrawlParam.FetchHint.Host
	}
	return doc.CrawlParam.FakeHost
}

