package maskyLogger

import (
	"errors"
	"fmt"
	"os"
)

// 日志类
type Logger struct {
	*LoggerConf
	callDepth int // 不放在LoggerConf中，因为同一个名称的Logger实例，Conf相同，可能被封装的层次不同，用CallDepth的不同
}

func (l *Logger) SetCallDepth(callDepth int) {
	l.callDepth = callDepth
}
func (l *Logger) GetCallDepth() int {
	if l.callDepth > 0 {
		return l.callDepth
	}

	return DefaultLoggerCallDepth
}
func (l *Logger) IsAll() bool    { return l.IsLogLevel(ALL) }
func (l *Logger) IsInfo() bool   { return l.IsLogLevel(LOG) }
func (l *Logger) IsDebug() bool  { return l.IsLogLevel(DEBUG) }
func (l *Logger) IsWarn() bool   { return l.IsLogLevel(WARN) }
func (l *Logger) IsError() bool  { return l.IsLogLevel(ERROR) }

func (l *Logger) IsLogLevel(level int) bool {
	return l.Levels[level]
}
func (l *Logger) checkAndLog(level int, args ...interface{}) {
	if l.Levels[level] {
		callDepth := l.GetCallDepth()
		levelStr := logLevelStringMap[level]
		for _, appender := range l.Appenders {
			appender.Log(callDepth, levelStr, args...)
		}
	}
}
func (l *Logger) checkAndLogf(level int, format string, args ...interface{}) {
	if l.Levels[level] {
		callDepth := l.GetCallDepth()
		levelStr := logLevelStringMap[level]
		for _, appender := range l.Appenders {
			appender.Logf(callDepth, levelStr, format, args...)
		}
	}
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.checkAndLogf(DEBUG, format, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.checkAndLogf(LOG, format, args...)
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.checkAndLogf(WARN, format, args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.checkAndLogf(ERROR, format, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.checkAndLogf(ERROR, format, args...)

	os.Exit(-1)
}

//日志的管理类，用于加载配置，返回正确的logger
type LoggerManager struct {
	*Tree
	cfg *Config
}

func NewLoggerManager(root *LoggerConf) *LoggerManager {
	manager := &LoggerManager{
		Tree: NewTree(root),
	}
	return manager
}
func NewLoggerManagerWithConf(cfg *Config) (*LoggerManager, error) {
	err := cfg.Verify()
	if err != nil {
		return nil, err
	}
	manager := NewLoggerManager(cfg.RootLogger())
	manager.cfg = cfg
	for _, loggerCfg := range cfg.BuildLoggers() {
		manager.insert(loggerCfg)
	}
	return manager, nil
}
func NewLoggerManagerWithJsconf(jsConf string) (*LoggerManager, error) {
	cfg, err := LoadConf(jsConf)
	if err != nil {
		return nil, err
	}
	return NewLoggerManagerWithConf(cfg)
}

func (lm *LoggerManager) UpdateConf(cfg *Config) error {
	err := cfg.Verify()
	if err != nil {
		return err
	}
	newTree := lm.Tree.clone()
	newTree.updateConf(cfg.RootLogger())
	for _, loggerCfg := range cfg.BuildLoggers() {
		newTree.updateConf(loggerCfg)
	}
	lm.cfg = cfg
	lm.Tree = newTree
	return nil
}
func (lm *LoggerManager) Logger(name string) (logger *Logger) {
	cfg := lm.inheritConf(name)
	cfg.Name = name
	return &Logger{
		LoggerConf: cfg,
	}
}

func (lm *LoggerManager) SetLogger(logger *Logger) {
	lm.Tree.updateConf(logger.LoggerConf)
}
func (lm *LoggerManager) SetRootAppender(appenders ...Appender) {
	lm.Root.current.SetAppender(appenders...)
	lm.Root.resetFinalConf()
}
func (lm *LoggerManager) UseRoot(name string) (err error) {
	if lm.cfg == nil {
		return errors.New("LoggerManager 缺少配置")
	}
	if root, ok := lm.cfg.RootsLogger()[name]; ok {
		lm.Root.current.Appenders = root.Appenders
		lm.Root.current.Levels = root.Levels
		lm.Root.resetFinalConf()
	} else {
		return fmt.Errorf("LoggerManager 找不到 [%s] Root Logger", name)
	}
	return nil
}
func (lm *LoggerManager) SetRootLevel(l int) {
	lm.Root.current.SetLevel(l)
	lm.Root.resetFinalConf()
}
func (lm *LoggerManager) SetRootOnlyLevel(ls ...int) {
	lm.Root.current.SetOnlyLevels(ls...)
	lm.Root.resetFinalConf()
}
