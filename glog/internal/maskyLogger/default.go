package maskyLogger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var defaultManager *LoggerManager

func init() {
	var defaultJson = `
		{
		    "Appenders":{
		        "stdout":{
		            "Type":"console"
		        }
		    },
		    "Root":{
		        "Level":"INFO",
		        "Appenders":["stdout"]
		    }
		}
		`
	var err error
	defaultManager, err = NewLoggerManagerWithJsconf(defaultJson)

	if value, exist := LogLevelMap[strings.ToUpper(os.Getenv("LOG_LEVEL"))]; exist {
		defaultManager.SetRootLevel(value)
	}

	if err != nil {
		panic(err)
	}
}

func Init(jsconf string) (err error) {
	cfg, err := LoadConf(jsconf)
	if err != nil {
		return err
	}
	return defaultManager.UpdateConf(cfg)
}
func InitConf(cfg *Config) (err error) {
	return defaultManager.UpdateConf(cfg)
}

func CurLoggerMananger() (cfg *LoggerManager) {
	return defaultManager
}

func Get(name string) *Logger {
	return defaultManager.Logger(name)
}

func UseRoot(name string) error {
	return defaultManager.UseRoot(name)
}
func SetRootAppender(appenders ...Appender) {
	defaultManager.SetRootAppender(appenders...)
}

func SetRootSeparationAppender(fileName string) {
	SetRootAppender(NewLevelSeparationDailyAppender("root", fileName, DefaultKeepDay))
}
func SetRootFileAppender(fileName string) {
	SetRootAppender(NewFileAppender("root", fileName))
}

func SetRootStrLevel(lstr string) (err error) {
	if i, ok := LogLevelMap[strings.ToUpper(strings.TrimSpace(lstr))]; ok {
		SetRootLevel(i)
	} else {
		err = fmt.Errorf("no such log level [%s]", lstr)
	}
	return
}
func SetRootLevel(l int)         { defaultManager.SetRootLevel(l) }
func SetRootOnlyLevel(ls ...int) { defaultManager.SetRootOnlyLevel(ls...) }

func Debugf(format string, args ...interface{}) { defaultLogger().Debugf(format, args...) }
func Infof(format string, args ...interface{})  { defaultLogger().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { defaultLogger().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { defaultLogger().Errorf(format, args...) }

func IsAll() bool    { return defaultLogger().IsAll() }
func IsInfo() bool   { return defaultLogger().IsInfo() }
func IsDebug() bool  { return defaultLogger().IsDebug() }
func IsWarn() bool   { return defaultLogger().IsWarn() }
func IsError() bool  { return defaultLogger().IsError() }

func defaultLogger() (logger *Logger) {
	name := pathInGoPath(2)
	logger = defaultManager.Logger(name)
	logger.SetCallDepth(DefaultLoggerCallDepth + 1)
	return
}
func pathInGoPath(level int) (inGoPath string) {
	_, name, _, _ := runtime.Caller(level + 1)
	if arr := strings.Split(name, "src/"); len(arr) > 1 {
		inGoPath = arr[1]
	} else {
		inGoPath = strings.Trim(name, "/")
	}
	return
}
