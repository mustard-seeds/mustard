package page_analysis
import (
	"testing"
	"strings"
	"mustard/base/string_util"
	"io/ioutil"
	"mustard/internal/github.com/PuerkitoBio/goquery"
)
///////init///////
func loadDoc(s string) (string,string) {
	path := strings.Split(s,"/")
	raw_url := path[len(path)-1]
	raw_url = strings.Replace(raw_url, "___","://", 1)
	raw_url = strings.Replace(raw_url, "_","/", -1)
	dat, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}
	return raw_url, string(dat)
}

func TestParseReg(t *testing.T) {
	var parser HtmlParser;
	u,c := loadDoc("./testdata/http___bj.58.com_zufang_24725387601214x.shtml")
	parser.RegisterRegex(`baidulon:'(\d+\.\d+)'`, func(i int, r[]string) {
		if r[1] != "116.46097179084" {
			t.Error("Regex parse fail." + r[1])
		}
	})
	parser.Parse(u, c)
}
func TestParseSelector(t *testing.T) {
	var parser HtmlParser;
	u,c := loadDoc("./testdata/http___bj.58.com_zufang_24725387601214x.shtml")
	parser.RegisterSelector("span.pay-method", func(i int, s *goquery.Selection){
		str := string_util.Purify(s.Text(), "\n","\t"," ")
		if str != "押一付三" {
			t.Error("Selector parse fail." + str)
		}
	})
	parser.Parse(u, c)
}
func TestParseSelectorWithKeyWord(t *testing.T) {
	var parser HtmlParser;
	u,c := loadDoc("./testdata/http___bj.58.com_zufang_24725387601214x.shtml")
	parser.RegisterSelectorWithTextKeyWord("span.pl10", "更新时间", func(i int, s *goquery.Selection){
		str := string_util.Purify(s.Text(), "\n","\t"," ")
		if str != "更新时间：2016-01-18" {
			t.Error("Selector Keyword parse fail." + str)
		}
	})
	parser.Parse(u, c)
}
