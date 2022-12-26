# log

A leveled, rotated logger.  
Wrapped [go.uber.org/zap](https://github.com/uber-go/zap)
and [natefinch/lumberjack](https://github.com/natefinch/lumberjack).

## Customize logger

This package's top levels log funcs use a global logger that initialized
by environment vars or can be configured at runtime with func
SetGlobalLoggerConfig.

Default config logs to stdout, log both levels info and debug. 
If logging to a file (set env LOG_FILE_PATH), the default config will rotate at
at midnight in UTC (7AM in Vietnam) or when its size reaches 100MB. Old log
files will be delete after 32 days 

## Usage

````go
package main

import "github.com/mywrap/log"

func main() {
	log.Debugf("line level debug 1+1 = %v", 1+1)
	log.Infof("line level info")
	name := "Dao Thanh Tung"
	log.Printf("hello %v", name)
}
````
Detail in [example.go](./example/example.go) and [log_test.go](./log_test.go).