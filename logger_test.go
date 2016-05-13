package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	log := NewLog(func(line *LogLine) {
		if line.Message != "test test" {
			t.Error(line.Message)
			t.Fail()
		}
	}, &Config{GoRoutinesLogTicker: 0 * time.Second})
	defer log.Close()

	log.Error("test %v", "test")
	log.Info("test test")
	log.Debug("%v", "test test")
	log.Notice("test %v%v%v%v", "t", "e", "s", "t")
}
