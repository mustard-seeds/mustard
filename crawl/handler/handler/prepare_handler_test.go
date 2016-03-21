package handler

import (
	"mustard/base/file"
	"mustard/crawl/proto"
	"testing"
)

func TestTranslateEncoding(t *testing.T) {
	content, _ := file.ReadFileToString("./testdata/gbk.html")
	doc := &proto.CrawlDoc{
		Content:     content,
		ContentType: "text/html",
	}
	translateEncoding(doc)
	if doc.OrigEncoding != "gbk" {
		t.Error("gbk html detect error.")
	}
	if doc.ConvEncoding != "utf-8" {
		t.Error("decode charset error..")
	}
}
