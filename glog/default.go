package glog

func init() {
	defaultLogger.SetCallDepth(defaultLogger.GetCallDepth() + 1)
}

// 默认logger
var defaultLogger = NewLogger("default")

// 默认 Debugf 打印方法
func Debugf(format string, args ...interface{}) { defaultLogger.Debugf(format, args...) }

// 默认 Infof 打印方法
func Infof(format string, args ...interface{}) { defaultLogger.Infof(format, args...) }

// 默认 Warnf 打印方法
func Warnf(format string, args ...interface{}) { defaultLogger.Warnf(format, args...) }

// 默认 Errorf 打印方法
func Errorf(format string, args ...interface{}) { defaultLogger.Errorf(format, args...) }

// 默认 Fatalf 打印方法
func Fatalf(format string, args ...interface{}) { defaultLogger.Fatalf(format, args...) }

// 默认 IsAll 方法
func IsAll() bool { return defaultLogger.IsAll() }

// 默认 IsInfo 方法
func IsInfo() bool { return defaultLogger.IsInfo() }

// 默认 IsDebug 方法
func IsDebug() bool { return defaultLogger.IsDebug() }

// 默认 IsWarn 方法
func IsWarn() bool { return defaultLogger.IsWarn() }

// 默认 IsError 方法
func IsError() bool { return defaultLogger.IsError() }
