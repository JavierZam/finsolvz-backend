package log

import (
	"context"
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "WARN: ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
)

func Debugf(ctx context.Context, format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	warnLogger.Printf(format, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Debug(ctx context.Context, msg string) {
	debugLogger.Println(msg)
}

func Info(ctx context.Context, msg string) {
	infoLogger.Println(msg)
}

func Warn(ctx context.Context, msg string) {
	warnLogger.Println(msg)
}

func Error(ctx context.Context, msg string) {
	errorLogger.Println(msg)
}

func Fatal(ctx context.Context, msg string) {
	errorLogger.Fatalln(msg)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	errorLogger.Fatalf(format, v...)
}
