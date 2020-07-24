package main

import (
"time"

"github.com/mywrap/log"
)

func main() {
	defer log.Flush()

	name := "Dao Thanh Tung"
	log.Printf("hello %v", name)

	log.SetGlobalLoggerConfig(log.Config{LogFilePath: "./log_test_file_example",
		RotateInterval: 24 * time.Hour, RotateRemainder: 7 * time.Hour})
	log.Debugf("line level debug 1+1 = %v", 1+1)
	log.Infof("line level info")
}
