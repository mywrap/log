package log

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// below tests only FAIL on panic

func TestStdLogger_Print(t *testing.T) {
	stdLogger := StdLogger{}
	stdLogger.Println("hihi", "ಠ_ಠ", "α,β,γ")
}

func TestPrintf(t *testing.T) {
	Printf("pussy")
}

func TestSetGlobalLoggerConfig(t *testing.T) {
	logFilePath0 := "./log_test_file.txt"
	os.Remove(logFilePath0)
	SetGlobalLoggerConfig(Config{LogFilePath: logFilePath0,
		RotateInterval: 24 * time.Hour, RotateRemainder: 7 * time.Hour})
	Info("test SetGlobalLoggerConfig")
}

func TestConcurrentlyLog(t *testing.T) {
	logFilePath1 := "./log_test_file_goroutine.txt"
	os.Remove(logFilePath1)
	SetGlobalLoggerConfig(Config{LogFilePath: logFilePath1,
		RotateInterval: 24 * time.Hour, RotateRemainder: 7 * time.Hour})
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			for k := 0; k < 100; k++ {
				Infof("test concurrently log i: %v, k: %v", i, k)
			}
		}()
	}
	wg.Wait()

	logData, err := ioutil.ReadFile(logFilePath1)
	if err != nil {
		t.Error(err)
	}
	nLines := strings.Count(string(logData), "\n")
	if nLines != 10000 {
		t.Errorf("nLines: expected: %v, real: %v", 10000, nLines)
	}
}

func TestTimedRotatingWriter(t *testing.T) {
	logFilePath2 := "./log_test_file_rotate"
	globPath2 := logFilePath2 + "*"
	countFileWildCard(globPath2, true)
	SetGlobalLoggerConfig(Config{
		LogFilePath:     logFilePath2,
		RotateInterval:  200 * time.Millisecond,
		RotateRemainder: 50 * time.Millisecond,
	})
	for k := 0; k < 60; k++ {
		time.Sleep(time.Duration(7+rand.Intn(6)) * time.Millisecond)
		Infof("test rotate log: k: %v", k)
	}
	globalLogger.rotator.close()
	nFiles, err := countFileWildCard(globPath2, false)
	if err != nil {
		t.Error(err)
	}
	if !(3 <= nFiles && nFiles <= 5) { // just estimate
		t.Errorf("rotate expected: %v, real: %v", 3, nFiles)
	}
}

func countFileWildCard(globPath string, delete bool) (int, error) {
	files, err := filepath.Glob(globPath)
	if err != nil {
		return 0, fmt.Errorf("err filepath Glob: %v", err)
	}
	nFiles := len(files)
	if !delete {
		return nFiles, nil
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return nFiles, fmt.Errorf("os Remove: %v", err)
		}
	}
	return nFiles, nil
}

func TestSetGlobalLoggerConfigWhileLogging(t *testing.T) {
	// TODO: TestSetGlobalLoggerConfigWhileLogging
}

func TestTruncateTime(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2021-04-12T09:29:00+07:00")
	lastRotated := calcLastRotatedTime(now, 24*time.Hour, 0)
	if r, e := lastRotated.Format(time.RFC3339), "2021-04-12T00:00:00Z"; r != e {
		t.Errorf("error calcLastRotatedTime: real %v, expected: %v", r, e)
	}
	vnLoc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Fatalf("error load vnLoc: %v", err)
	}
	now2, _ := time.Parse(time.RFC3339, "2021-04-12T02:29:00Z")
	lastRotated2 := calcLastRotatedTime(now2, 24*time.Hour, 17*time.Hour).In(vnLoc)
	if r, e := lastRotated2.Format(time.RFC3339), "2021-04-12T00:00:00+07:00"; r != e {
		t.Errorf("error calcLastRotatedTime: real %v, expected: %v", r, e)
	}
}
