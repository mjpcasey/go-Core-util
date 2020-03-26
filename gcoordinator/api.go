package gcoordinator

import (
	"gcore/app/share"
	"gcore/gcoordinator/gcoordinatorTypes"
	"gcore/gcoordinator/internal/zkCoordinator"
	"gcore/glog"
	"net/url"
	"strings"
)

var logger = glog.NewLogger(`Coordinator`)

// NewCoordinator create Coordinator object base on app config
func NewCoordinator(app *share.AppShare, confPrefix string) (c Coordinator) {
	uri, err := url.ParseRequestURI(app.Conf.Get(confPrefix + "url"))

	if err == nil {
		cfg := &gcoordinatorTypes.Config{
			Root:  app.Conf.Get(confPrefix + "root"),
			Addrs: strings.Split(uri.Host, ","),
		}
		switch uri.Scheme {
		case "zookeeper":
			c = zkCoordinator.New(cfg)
		default:
			logger.Fatalf("未知的协调器类型: %s", uri.Scheme)
		}
	} else {
		logger.Fatalf("解析协调器配置[%s]出错: %s", app.Conf.Get(confPrefix+"url"), err)
	}

	return
}
