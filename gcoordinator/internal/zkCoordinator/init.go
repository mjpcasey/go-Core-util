package zkCoordinator

import (
	"gcore/glog"
	"github.com/samuel/go-zookeeper/zk"
)

var logger = glog.NewLogger("zookeeperCoordinator")

type zkLogger struct{}

func (l *zkLogger) Printf(fmt string, vals ...interface{}) {
	//logger.Errorf(fmt, vals...)
}

func init() {
	// 替换三方包 logger
	zk.DefaultLogger = new(zkLogger)
}