package logger

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"sync"
	"time"
)

// Log - logger struct
type Log struct {
	LogStream   chan *LogLine
	processor   func(*LogLine)
	wg          *sync.WaitGroup
	LogSeverity map[Severity]bool
}

// Severity - type for severity
type Severity string

const (
	// ERROR - error severity
	ERROR = Severity("error")
	// WARNING - warning severity
	WARNING = Severity("warning")
	// INFO - info severity
	INFO = Severity("info")
	// NOTICE - notice severity
	NOTICE = Severity("notice")
	// DEBUG - debug severity
	DEBUG = Severity("debug")
)

var colorMap = map[Severity]func(...interface{}) string{
	ERROR:   color.New(color.FgRed).SprintFunc(),
	WARNING: color.New(color.FgYellow).SprintFunc(),
	INFO:    color.New(color.FgGreen).SprintFunc(),
	NOTICE:  color.New(color.FgCyan).SprintFunc(),
	DEBUG:   color.New(color.FgBlue).SprintFunc(),
}

// ColoredString - will output severity as a nice colorfull string
func (s *Severity) ColoredString() string {
	return colorMap[*s](string(*s))
}

// LogLine - struct containing info about log
type LogLine struct {
	Message  string
	Severity Severity
}

// Print - will print logline to stdout
func (l *LogLine) Print() {
	fmt.Printf("[%v] %v\n", l.Severity.ColoredString(), l.Message)
}

// NewLog - creates new logger
func NewLog(processor func(line *LogLine)) *Log {
	stream := make(chan *LogLine)
	wg := &sync.WaitGroup{}

	l := &Log{
		LogStream:   stream,
		LogSeverity: map[Severity]bool{INFO: true, ERROR: true, WARNING: true, NOTICE: true, DEBUG: false},
		processor:   processor,
		wg:          wg,
	}

	go func(l *Log) {
		for {
			ticker := time.NewTicker(30 * time.Second)
			select {
			case <-ticker.C:
				l.Debug(fmt.Sprintf("Goroutines count: %v", runtime.NumGoroutine()))
			}
		}
	}(l)

	wg.Add(1)
	go func(stream chan *LogLine, processor func(line *LogLine), wg *sync.WaitGroup) {
		defer wg.Done()
		for line := range stream {
			processor(line)
		}
	}(stream, processor, wg)

	return l
}

func (l *Log) log(severity Severity, m string) {
	if l.LogSeverity[severity] {
		l.LogStream <- &LogLine{m, severity}
	}
}

// Notice - puts notice into chan
func (l *Log) Notice(m string) {
	l.log(NOTICE, m)
}

// Error - puts error into chan
func (l *Log) Error(m string) {
	l.log(ERROR, m)
}

// Info - puts info into chan
func (l *Log) Info(m string) {
	l.log(INFO, m)
}

// Warning - puts warning into chan
func (l *Log) Warning(m string) {
	l.log(WARNING, m)
}

// Debug - puts debug into chan
func (l *Log) Debug(m string) {
	l.log(DEBUG, m)
}

// Close - will close log and wait for log processor to finish
func (l *Log) Close() {
	close(l.LogStream)
	l.wg.Wait()
}
