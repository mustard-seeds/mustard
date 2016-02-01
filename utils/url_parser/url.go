package url_parser

import (
	"net/url"
	"mustard/internal/github.com/PuerkitoBio/purell"
	LOG "mustard/base/log"
)

func GetHost(_url string) string {
	var u *url.URL
	var e error
	if u, e = url.Parse(_url); e != nil {
		LOG.Error("Parse Url Fail")
		return nil
	}
	return u.Host
}
func NormalizeUrl(_url string) string {
	normal_url, err := purell.NormalizeURLString(_url, purell.FlagsSafe)
	if err != nil {
		return nil
	}
	return normal_url
}
