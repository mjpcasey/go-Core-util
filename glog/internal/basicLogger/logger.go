package basicLogger

import (
	"fmt"
	"log"
	"os"
)

var (
	defaultCallDepth = 2
)

type logger struct {
	name string

	stdout *log.Logger
	stderr *log.Logger

	callDepth int
}

func (lg *logger) format(level LevelType, msg string, args []interface{}) string {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return "[" + level.String() + "] " + msg
}

func (lg *logger) IsAll() bool {
	return threshold <= LevelDebug
}
func (lg *logger) IsDebug() bool {
	return threshold <= LevelDebug
}
func (lg *logger) IsInfo() bool {
	return threshold <= LevelInfo
}
func (lg *logger) IsWarn() bool {
	return threshold <= LevelWarn
}
func (lg *logger) IsError() bool {
	return threshold <= LevelError
}
func (lg *logger) IsFatal() bool {
	return threshold <= LevelFatal
}

func (lg *logger) Debugf(fmt string, args ...interface{}) {
	if lg.IsDebug() {
		_ = lg.stdout.Output(lg.GetCallDepth(), lg.format(LevelDebug, fmt, args))
	}
}
func (lg *logger) Infof(fmt string, args ...interface{}) {
	if lg.IsInfo() {
		_ = lg.stdout.Output(lg.GetCallDepth(), lg.format(LevelInfo, fmt, args))
	}
}
func (lg *logger) Warnf(fmt string, args ...interface{}) {
	if lg.IsWarn() {
		_ = lg.stderr.Output(lg.GetCallDepth(), lg.format(LevelWarn, fmt, args))
	}
}
func (lg *logger) Errorf(fmt string, args ...interface{}) {
	if lg.IsError() {
		_ = lg.stderr.Output(lg.GetCallDepth(), lg.format(LevelError, fmt, args))
	}
}
func (lg *logger) Fatalf(fmt string, args ...interface{}) {
	if lg.IsFatal() {
		_ = lg.stderr.Output(lg.GetCallDepth(), lg.format(LevelFatal, fmt, args))
	}

	os.Exit(-1)
}

func (lg *logger) SetCallDepth(depth int) {
	lg.callDepth = depth
}
func (lg *logger) GetCallDepth() int {
	if lg.callDepth > 0 {
		return lg.callDepth
	}

	return defaultCallDepth
}

func New(name string) (lg *logger) {
	lg = &logger{
		name:   name,
		stdout: log.New(os.Stdout, "", log.LstdFlags|log.Llongfile),
		stderr: log.New(os.Stderr, "", log.LstdFlags|log.Llongfile),
	}

	return
}
