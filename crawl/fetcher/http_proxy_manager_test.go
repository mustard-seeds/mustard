package fetcher

import (
	"testing"
)

func getrrHttpProxyManager() *ProxyManager {
	*CONF.Crawler.ProxyConfFile = "../../mdata/etc/crawl/fetch_proxys.config"
	return NewProxyManager(PROXY_SELECT_RR)
}

func TestRR(t *testing.T) {
	m := getrrHttpProxyManager()
	if len(m.alives) == 0 {
		t.Error("Config file does not have proxy...")
	}
	if len(m.deads) != 0 {
		t.Error("Init should be all alive")
	}
	p, e := m.GetProxyUrl()
	if e != nil {
		t.Errorf("get proxy url error.%s", e.Error())
	}
	if len(m.alives) > 1 {
		p1, e2 := m.GetProxyUrl()
		if e2 != nil {
			t.Errorf("get second proxy url error %s", e.Error())
		}
		if p.Host == p1.Host {
			t.Error("get two same proxy?")
		}
	}
	m.MarkDeadProxy(p)
	if len(m.deads) != 1 {
		t.Error("Mark dead proxy error.")
	}
}
