# log

A leveled, rotated logger.  
Wrapped [go.uber.org/zap](https://github.com/uber-go/zap)
and [natefinch/lumberjack](https://github.com/natefinch/lumberjack).

## Customize logger

This package's top levels log funcs use a global logger that initialized
by environment vars or can be configured at runtime with func
SetGlobalLoggerConfig.

Default logger logs to stdout, log both levels are info and debug. 
If logging to a file, the default logger will rotate at midnight (+07:00)
or when its size reaches 100MB.  

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