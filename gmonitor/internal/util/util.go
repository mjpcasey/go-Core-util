package util

import (
	"gcore/glog"
	"regexp"
	"sync"
)

var metricMap sync.Map

func isLegal(name string) bool {
	return regexp.MustCompile(`^[a-z0-9A-Z_]+$`).MatchString(name)
}

func Check(name string) {
	if isLegal(name) != true {
		glog.Fatalf("监控指标 %s 命名不合规则,应只包含大小写字母_数字和下划线", name)
	}

	if _, loaded := metricMap.LoadOrStore(name, true); loaded == true {
		glog.Fatalf("监控指标\"%s\"已经存在", name)
	}
}
