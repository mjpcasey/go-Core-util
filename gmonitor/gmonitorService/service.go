package gmonitorService

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gcore/app/share"
	"gcore/ghttp"
	"gcore/glog"
)

var logger = glog.NewLogger(`monitorService`)

// Service gcore监控
type Service struct {
	confPath   string
	appShare   *share.AppShare
	httpServer *http.Server
}

// startHTTP 初始化 prometheus 的http监控接口
func (m *Service) startHTTP(port int) {
	router := httprouter.New()

	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/service/healthcheck", handlerHealthCheck)

	m.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ghttp.NewHandler(router),
	}

	go func() {
		if err := m.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("指标监控器接口启动错误: %s", err)
		}
	}()

	logger.Debugf("启动指标监测数据接口[端口=%d]", port)
}

func (m *Service) stopHTTP() (err error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second)

	return m.httpServer.Shutdown(ctx)
}

// Stop 关闭接口
func (m *Service) Stop() {
	err := m.stopHTTP()

	if err != nil {
		logger.Errorf("关闭指标监测接口出错: %s", nil)
	}
}

// Start 开启接口
func (m *Service) Start() {
	port := m.appShare.Conf.GetIntDef(m.confPath+"httpPort", 9500)
	m.startHTTP(port)
}

// NewService 初始化监控服务
func NewService(share *share.AppShare, configPrefix string) (m *Service) {
	m = &Service{
		appShare: share,
		confPath: configPrefix,
	}

	return m
}
