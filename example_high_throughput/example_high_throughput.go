package main

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/mywrap/log"
)

// On Windows, reading a file while it is being written by another process
// may cause errors. This program continuously write a large amount of data
// to file "log_test_file.txt", you can try to read the file.

// require: "set LOG_FILE_PATH=log_test_file.txt"

func main() {
	log.Printf("generating random data")
	const nRoutines = 1000
	const nLines = 25
	data := make([]string, nRoutines*nLines)
	for i := 0; i < nRoutines*nLines; i++ {
		data[i] = GenRandomWord(10, 20000, AlphaEnList)
		// 10 KB/line * 25000 line so all data is about 250MB
	}

	wg := &sync.WaitGroup{}
	beginT := time.Now()
	log.Printf("begin writing log")
	for r := 0; r < nRoutines; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Add(-1)
			for i := 0; i < nLines; i++ {
				line := data[r*nLines+i]
				log.Printf("i %9d: %v", i, line)
			}
		}(r)
	}
	wg.Wait()
	log.Printf("end writing log. duration: %v", time.Since(beginT).Seconds())
}

var AlphaEnList = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenRandomWord(minLen int, maxLen int, charList []rune) string {
	if minLen <= 0 {
		minLen = 0
	}
	if maxLen < minLen {
		maxLen = minLen
	}
	wordLen := minLen + rand.Intn(maxLen+1-minLen)
	builder := strings.Builder{}
	builder.Grow(3 * wordLen) // UTF8
	for i := 0; i < wordLen; i++ {
		builder.WriteRune(charList[rand.Intn(len(charList))])
	}
	return builder.String()
}
