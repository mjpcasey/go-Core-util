package basicLogger

import (
	"os"
	"strings"
)

type LevelType int

func (level LevelType) String() string {
	return levelTexts[level]
}

const (
	LevelDebug LevelType = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelTexts = map[LevelType]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "LOG",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

var threshold = LevelInfo

func SetLogLevel(level string) {
	for value, text := range levelTexts {
		if strings.EqualFold(level, text) {
			threshold = value
			return
		}
	}
}

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))
}
