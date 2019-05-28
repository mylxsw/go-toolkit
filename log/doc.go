/*
Package log 实现了不同级别的日志输出，可以针对不通的模块单独定制日志输出格式。

	package main

	import(
		"github.com/mylxsw/go-toolkit/log"
	)

	var logger = log.Module("toolkit.process")

	func main() {
		logger.Debugf("xxxx: %s, xxx", "ooo")
		logger.WithContext(log.C{
			"id": 123,
			"name": "lixiaoyao",
		}).Debugf("Hello, %s", "world")
	}

*/
package log
