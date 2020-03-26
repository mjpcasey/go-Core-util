/*
Package glog is a basic logger library for simple golang projects.
*/
package glog

// Logger is the interface of logger instances.
type Logger interface {
	IsAll() bool                            // IsAll 日志等级判断
	IsDebug() bool                          // IsDebug 日志等级判断
	IsInfo() bool                           // IsInfo 日志等级判断
	IsWarn() bool                           // IsWarn 日志等级判断
	IsError() bool                          // IsError 日志等级判断
	Debugf(format string, a ...interface{}) // Debug 信息打印接口
	Infof(format string, a ...interface{})  // Info 信息打印接口
	Warnf(format string, a ...interface{})  // Warn 信息打印接口
	Errorf(format string, a ...interface{}) // Error 信息打印接口
	Fatalf(format string, a ...interface{}) // Fatal 信息打印接口
	SetCallDepth(callDepth int)
	GetCallDepth() int
}
