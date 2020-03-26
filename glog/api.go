package glog

import (
	"gcore/glog/internal/basicLogger"
	"gcore/glog/internal/maskyLogger"
)

const (
	DEBUG = "DEBUG"
	INFO  = "LOG"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

func SetLogLevel(level string) {
	basicLogger.SetLogLevel(level)
	maskyLogger.SetRootStrLevel(level)
}

// NewLogger creates a new logger instance match Logger interface
func NewLogger(name string) Logger {
	if name != "" {
		name = "[" + name + "] "
	}

	return basicLogger.New(name)
	//return maskyLogger.Get(name)
}