package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type IUserAgentGenerator interface {
	Get() string
}

type UserAgentGenerator struct {
}

func (u *UserAgentGenerator) Get() string {
	return "Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25"
}

type Logger struct {
	ch  chan string
	dst io.WriteCloser
}

func (l *Logger) Log(s string) {
	l.ch <- s
}

func (l *Logger) Close() {
	close(l.ch)
	l.dst.Close()
}

func (l *Logger) loop() {
	for {
		s, ok := <-l.ch
		if !ok {
			return
		}
		l.dst.Write([]byte(fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05 -07:00"), s)))
	}
}

func NewFileLogger() *Logger {
	file, err := os.OpenFile(fmt.Sprintf("log_%s.txt", time.Now().Format("20060102")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l := &Logger{make(chan string), file}
	go l.loop()
	return l
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_."

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
