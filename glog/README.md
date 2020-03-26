日志包

用法：

使用环境变量控制日志等级

LOG_LEVEL=debug|info|warn|error|fatal

```go

var logger = glog.NewLogger("name")

logger.Debugf("for %s", "Debug")
logger.Infof("for %s", "Info")
logger.Warnf("for %s", "Warn")
logger.Errorf("for %s", "Error")
logger.Fatalf("for %s", "Fatal")
```

提高性能：
```go
if logger.IsDebug() {
	logger.Debugf("for %s", "Debug")
}
if logger.IsInfo() {
	logger.Infof("for %s", "Debug")
}
if logger.IsWarn() {
	logger.Warnf("for %s", "Debug")
}
if logger.IsError() {
	logger.Errorf("for %s", "Debug")
}
```