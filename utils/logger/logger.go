package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type ILogger interface {
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type LoggerOptions struct {
	DateFormat *string
}

type FileLogger struct {
	opts    LoggerOptions
	mux     sync.Mutex
	info    *os.File
	warning *os.File
	error   *os.File
	fatal   *os.File
}

func (l *FileLogger) output(dst *os.File, s string) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.opts.DateFormat != nil {
		dst.Write([]byte(fmt.Sprintf("[%s] %s\n", time.Now().Format(*l.opts.DateFormat), s)))
	} else {
		dst.Write([]byte(fmt.Sprintf("%s\n", s)))
	}
}
func (l *FileLogger) Info(args ...interface{}) {
	l.output(l.info, "INFO: "+fmt.Sprint(args...))
}
func (l *FileLogger) Warning(args ...interface{}) {
	l.output(l.info, "WARNING: "+fmt.Sprint(args...))
}
func (l *FileLogger) Error(args ...interface{}) {
	l.output(l.info, "ERROR: "+fmt.Sprint(args...))
}
func (l *FileLogger) Fatal(args ...interface{}) {
	l.output(l.info, "FATAL: "+fmt.Sprint(args...))
	os.Exit(2)
}
func (l *FileLogger) Infof(f string, args ...interface{}) {
	l.output(l.info, fmt.Sprintf("INFO: "+f, args...))
}
func (l *FileLogger) Warningf(f string, args ...interface{}) {
	l.output(l.info, fmt.Sprintf("WARNING: "+f, args...))
}
func (l *FileLogger) Errorf(f string, args ...interface{}) {
	l.output(l.info, fmt.Sprintf("ERROR: "+f, args...))
}
func (l *FileLogger) Fatalf(f string, args ...interface{}) {
	l.output(l.info, fmt.Sprintf("FATAL: "+f, args...))
	os.Exit(2)
}

func mergeOptions(opts ...LoggerOptions) LoggerOptions {
	res := LoggerOptions{}
	for _, o := range opts {
		res.DateFormat = o.DateFormat
	}
	return res
}

func NewFileLogger(o ...LoggerOptions) ILogger {
	file, err := os.OpenFile(fmt.Sprintf("log_%s.txt", time.Now().Format("20060102")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l := &FileLogger{}
	l.mux = sync.Mutex{}
	l.info = file
	l.error = file
	l.warning = file
	l.fatal = file
	l.opts = mergeOptions(o...)
	return l
}

var defaultTimeFormat = "2006-01-02 15:04:05 -07:00"
var defaultLogger = NewFileLogger(LoggerOptions{&defaultTimeFormat})

func DefaultLogger() ILogger {
	return defaultLogger
}
