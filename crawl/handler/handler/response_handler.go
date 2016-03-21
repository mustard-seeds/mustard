/*
Send Crawldoc to Crawldoc.CrawlParam.Receivers
*/
package handler

import (
	"mustard/crawl/proto"
	"mustard/crawl/scheduler"
	"mustard/base/string_util"
	"mustard/base/time_util"
	"mustard/base/hash"
)

const (
	kCrawlDocSenderCouldRetryInterval int64 = 2
	kDefaultResponseHandlerSendSpeed = 1000
)

type ResponseHandler struct {
	CrawlHandler
	// client cache, could reconnect.
	connectionPools map[string]*scheduler.CrawlDocSender
}

func (handler *ResponseHandler) Init() bool {
	//TODO response handler, add crawldoc sender
	return false
}
func (handler *ResponseHandler) Accept(crawlDoc *proto.CrawlDoc) bool {
	return crawlDoc.CrawlParam != nil && len(crawlDoc.CrawlParam.Receivers) > 0
}
func (handler *ResponseHandler) Process(crawlDoc *proto.CrawlDoc) {
	for true {
		sendSuccess := false
		// random select receiver.
		for _, idx := range hash.ShuffleInt(len(crawlDoc.CrawlParam.Receivers)) {
			conn := crawlDoc.CrawlParam.Receivers[idx]
			var serverAddr string
			string_util.StringAppendF(&serverAddr, "%s:%d", conn.Host, conn.Port)
			sender, exist := handler.connectionPools[serverAddr]
			if !exist {
				sender = scheduler.NewCrawlDocSender(
					conn.Host,
					int(conn.Port),
					kDefaultResponseHandlerSendSpeed)
				sender.Init()
				handler.connectionPools[serverAddr] = sender
			}
			if sender.Connected ||
			time_util.GetCurrentTimeStamp() - sender.LastReconnectTimeStamp > kCrawlDocSenderCouldRetryInterval {
				sender.Flush(crawlDoc)
				sendSuccess = true
				break
			}
		}
		if sendSuccess {
			break
		}
		time_util.Sleep(1)
	}
}

// use for create instance from a string
func init() {
	registerCrawlTaskType(&ResponseHandler{})
}
