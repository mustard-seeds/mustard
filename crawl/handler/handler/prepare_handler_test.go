package handler

import (
	"testing"
	"mustard/base/file"
	"mustard/crawl/proto"
	"mustard/internal/golang.org/x/net/html/charset"
)

func TesttranslateEncoding(t *testing.T) {
	content,_ := file.ReadFileToString("./testdata/gbk.html")
	doc := &proto.CrawlDoc{
		Content:content,
		ContentType:"text/html",
	}
	translateEncoding(doc)
	if doc.OrigEncoding != "gbk" {
		t.Error("gbk html detect error.")
	}
	if doc.ConvEncoding != "utf-8" {
		t.Error("decode charset error..")
	}

	_, n, _ := charset.DetermineEncoding(doc.Content, doc.ContentType)
	if n != "utf-8" {
		t.Errorf("decode charset not to utf-8 is %s", n)
	}
}