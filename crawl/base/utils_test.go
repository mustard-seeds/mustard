package base

import (
	"mustard/crawl/proto"
	"testing"
)

func TestGetHostName(t *testing.T) {
	doc := &proto.CrawlDoc{
		CrawlParam: &proto.CrawlParam{
			FetchHint: &proto.FetchHint{
				Host: "x.com",
			},
		},
	}
	h := GetHostName(doc)
	if h != "x.com" {
		t.Error("GetHostName without FakeHost not work")
	}
	doc.CrawlParam.FakeHost = "a.com"
	hh := GetHostName(doc)
	if hh != "a.com" {
		t.Error("GetHostName with FakeHost not work")
	}
}

func TestGetDomainFromHost(t *testing.T) {
	pair := map[string]string{
		"www.sina.com":         "sina.com",
		"z.cn":                 "z.cn",
		"bj.58.com":            "58.com",
		"aaa.blog.sina.com.cn": "blog.sina.com.cn",
	}
	for k, v := range pair {
		if v != GetDomainFromHost(k) {
			t.Errorf("GetDomainFromHost Error, %s -> %s", k, v)
		}
	}
}

func TestIsSameDomainGood(t *testing.T) {
	pair := map[string]string{
		"cn":                 "cn",
		"xx.com":             "xx.com",
		"auto.blog.sina.com": "blog.sina.com",
		"xxx.blog.sina.com":  "sina.com",
	}
	for k, v := range pair {
		if !IsSameDomain(k, v) {
			t.Errorf("No same Domain? %s <-> %s", k, v)
		}
	}
}

func TestIsSameDomainBad(t *testing.T) {
	pair := map[string]string{
		"cn":                "z.cn",
		"xx.com":            "xxx.com",
		"xxx.blog.sina.com": "sina.com.cn",
	}
	for k, v := range pair {
		if IsSameDomain(k, v) {
			t.Errorf("Same Domain? %s <-> %s", k, v)
		}
	}
}
