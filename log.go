// Package log provides a leveled, rotated, fast, structured logger.
// This package APIs Print and Fatal are compatible the standard library log.
package log

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	rotator *timedRotatingWriter
}

// globalLogger will be inited env vars (detail in func newConfigFromEnv)
// All funcs in this package use globalLogger
var globalLogger *Logger = NewLogger(newConfigFromEnv())

// Config will be used for initializing a Logger
type Config struct {
	LogFilePath string // default log to stdout
	// default log both info and debug logs,
	// change this field to true will only log info.
	IsLogLevelInfo bool
	// whether to log simultaneously to both stdout and file
	IsNotLogBoth bool
	// default 24 hours (rotating once per day if size of the log file < 100MB)
	RotateInterval time.Duration
	// default rotate at midnight in UTC (or 7AM in Vietnam)
	RotateRemainder time.Duration
}

// newConfigFromEnv create a Config from env vars
func newConfigFromEnv() Config {
	var c Config
	c.LogFilePath = os.Getenv("LOG_FILE_PATH")
	c.IsLogLevelInfo, _ = strconv.ParseBool(os.Getenv("LOG_LEVEL_INFO"))
	c.IsNotLogBoth, _ = strconv.ParseBool(os.Getenv("LOG_NOT_STDOUT"))
	c.RotateInterval = 24 * time.Hour
	c.RotateRemainder = 0
	return c
}

// NewLogger returns a inited Logger
func NewLogger(conf Config) *Logger {
	ret := &Logger{}

	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConf)

	var writers []zapcore.WriteSyncer
	stdWriter, _, _ := zap.Open("stdout")
	if conf.LogFilePath == "" {
		writers = []zapcore.WriteSyncer{stdWriter}
	} else {
		ret.rotator = newTimedRotatingWriter(conf.LogFilePath,
			conf.RotateInterval, conf.RotateRemainder)
		fileWriter := zapcore.AddSync(ret.rotator)
		if conf.IsNotLogBoth {
			writers = []zapcore.WriteSyncer{fileWriter}
		} else { // default behavior
			writers = []zapcore.WriteSyncer{stdWriter, fileWriter}
		}
	}
	combinedWriter := zap.CombineWriteSyncers(writers...)

	logLevel := zap.DebugLevel
	if conf.IsLogLevelInfo {
		logLevel = zap.InfoLevel
	}

	core := zapcore.NewCore(encoder, combinedWriter, logLevel)
	zl := zap.New(core, zap.AddCaller())
	zl = zl.WithOptions(zap.AddCallerSkip(1))
	ret.SugaredLogger = zl.Sugar()
	return ret
}

// SetGlobalLoggerConfig replace the globalLogger with a new logger from conf
func SetGlobalLoggerConfig(conf Config) {
	globalLogger.SugaredLogger.Sync()
	if globalLogger.rotator != nil {
		globalLogger.rotator.close()
	}
	customizedLogger := NewLogger(conf)
	globalLogger = customizedLogger
}

// Flush buffered log lines of globalLogger
func Flush() {
	globalLogger.SugaredLogger.Sync()
}

type timedRotatingWriter struct {
	*lumberjack.Logger
	interval    time.Duration
	remainder   time.Duration
	mutex       sync.RWMutex
	lastRotated time.Time
}

func calcLastRotatedTime(now time.Time, interval, remainder time.Duration) time.Time {
	nowUTC := now.UTC()
	return nowUTC.Add(-remainder).Truncate(interval).Add(remainder)
}

func newTimedRotatingWriter(filePath string,
	interval time.Duration, remainder time.Duration) *timedRotatingWriter {
	base := &lumberjack.Logger{
		Filename: filePath,
		MaxSize:  100, // MaxSize unit is MiB
		MaxAge:   32,  // MaxAge unit is days
	}
	w := &timedRotatingWriter{Logger: base, interval: interval, remainder: remainder}
	w.mutex.Lock()
	w.Logger.Rotate()
	w.lastRotated = calcLastRotatedTime(time.Now(), w.interval, w.remainder)
	w.mutex.Unlock()
	return w
}

func (w *timedRotatingWriter) rotateIfNeeded() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if time.Since(w.lastRotated) < w.interval {
		return nil
	}
	w.lastRotated = calcLastRotatedTime(time.Now(), w.interval, w.remainder)
	err := w.Logger.Rotate()
	return err
}

// Write implements io.Writer interface,
func (w *timedRotatingWriter) Write(p []byte) (int, error) {
	err := w.rotateIfNeeded()
	if err != nil {
		return 0, err
	}
	// ensure no goroutine write log while rotating
	w.mutex.RLock()
	n, err := w.Logger.Write(p)
	w.mutex.RUnlock()
	return n, err
}

// close the log file
func (w *timedRotatingWriter) close() {
	w.Logger.Close()
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Print(args ...interface{}) {
	globalLogger.Info(args...)
}

func Println(args ...interface{}) {
	globalLogger.Info(args...)
}

func Printf(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Condf(cond bool, format string, args ...interface{}) {
	if cond {
		globalLogger.Infof(format, args...)
	}
}

// StdLogger is compatible with the standard library logger,
// This logger call the globalLogger funcs
type StdLogger struct{}

func padArgs(args []interface{}) []interface{} {
	if len(args) <= 1 {
		return args
	}
	newArgs := make([]interface{}, 2*len(args)-1)
	for i, e := range args {
		newArgs[2*i] = e
		if i != len(args)-1 {
			newArgs[2*i+1] = " "
		}
	}
	return newArgs
}

func (l StdLogger) Print(args ...interface{}) {
	globalLogger.Info(padArgs(args)...)
}

func (l StdLogger) Println(args ...interface{}) {
	globalLogger.Info(padArgs(args)...)
}

func (l StdLogger) Printf(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func (l *StdLogger) Fatal(v ...interface{}) {
	globalLogger.Fatal(padArgs(v)...)
}

func (l *StdLogger) Fatalln(v ...interface{}) {
	globalLogger.Fatal(padArgs(v)...)
}

func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

func (l *StdLogger) Panic(v ...interface{}) {
	globalLogger.Info(v...)
	panic(1)
}

func (l *StdLogger) Panicf(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
	panic(1)
}

func (l *StdLogger) Panicln(v ...interface{}) {
	globalLogger.Info(v...)
	panic(1)
}
