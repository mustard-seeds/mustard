package main

import (
	"mustard/base/conf"
	"mustard/base/file"
	LOG "mustard/base/log"
	"mustard/base/proto_util"
	pb "mustard/crawl/proto"
	"mustard/crawl/scheduler"
	"strings"
)

var CONF = conf.Conf

func main() {
	// crawldoc sender
	sender := scheduler.NewCrawlDocSender(*CONF.Crawler.DispatcherHost,
		*CONF.Crawler.DispatcherPort,
		*CONF.Crawler.DefaultSendSpeed)
	sender.Init()

	// params filler init
	filler := scheduler.ParamFillerMaster{}
	filler.RegisterJobDescription(scheduler.GetJobDescriptionFromFile(*CONF.Crawler.JobDescriptionConfFile))
	filler.RegisterParamFillerGroup(&scheduler.DefaultParamFillerGroup{})
	filler.Init()

	// callback from read file
	fname := file.GetConfFile(*CONF.Crawler.UrlScheduleFile)
	file.FileLineReader(fname, "#", func(line string) {
		if !strings.HasPrefix(line, "http") {
			LOG.VLog(1).Debugf("Error url format %s", line)
			return
		}
		doc := pb.CrawlDoc{RequestUrl: line}
		filler.Fill(&doc)
		sender.Flush(&doc)
		LOG.VLog(3).Debugf("Send one doc: %s", doc.RequestUrl)
		LOG.VLog(3).Debugf("Send one doc Detail:\n %s", proto_util.FromProtoToString(&doc))
	})
}
