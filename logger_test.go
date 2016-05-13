package logger

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	wg := sync.WaitGroup{}
	counter := 0
	log := NewLog(func(line *LogLine) {
		defer wg.Done()
		fmt.Println(line.File)
		if line.Message != "test test" {
			t.Error(line.Message)
			t.Fail()
		}
		counter++
	}, &Config{GoRoutinesLogTicker: 0 * time.Second})
	defer log.Close()

	wg.Add(4)
	log.Error("test %v", "test")
	log.Info("test test")
	log.Warning("%v", "test test")
	log.Notice("test %v%v%v%v", "t", "e", "s", "t")

	log.Debug("Shouldn't log at all")

	wg.Wait()

	if counter != 4 {
		t.Error(counter)
		t.Fail()
	}
}
