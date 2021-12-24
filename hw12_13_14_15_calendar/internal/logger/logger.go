package logger

import "fmt"

type Logger struct { // TODO
}

func New(level string) *Logger {
	return &Logger{}
}

func (l Logger) Info(msg ...interface{}) {
	fmt.Println(msg...)
}

func (l Logger) Infof(str string, msg ...interface{}) {
	fmt.Printf(str+"\n", msg...)
}

func (l Logger) Log(msg ...interface{}) {
	fmt.Println(msg...)
}

func (l Logger) Error(msg ...interface{}) {
	fmt.Println(msg...)
}

// TODO
