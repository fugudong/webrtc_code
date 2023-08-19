package main

import (
	log "github.com/cihub/seelog"
)

func SetupLogger() {
	logger, err := log.LoggerFromConfigAsFile("seelog2.xml")
	if err != nil {
		return
	}

	log.ReplaceLogger(logger)

}

func testFlush()  {
	defer log.Flush()
	log.Debug("一定要flush?")	// 确实是一定要flush
}
// seelog 使用笔记 https://www.liangzl.com/get-article-detail-16487.html
func main() {
	//defer log.Flush()
	SetupLogger()
	testFlush()
	log.Info("Hello from Seelog!")
	log.Debug("我是真心想用你")
}