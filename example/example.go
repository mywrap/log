package main

import (
	"github.com/mywrap/log"
)

func main() {
	// set environment var LOG_FILE_PATH to log to file

	name := "Dao Thanh Tung"
	log.Printf("hello %v", name)

	// set environment var LOG_LEVEL_INFO=true to skip following log line
	log.Debugf("hi, this is a debug level log line")

	log.Printf("done")
}
