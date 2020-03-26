/**
Package app 用于提供 gcore 框架的初始化，并提供应用相关方法

提供 GetDBManager 用于获取数据库管理器

提供 GetRpcManager 用于获取RPC管理器

提供 OnStart OnStop 用于设置业务线程的启动关闭
*/
package app

import (
	"gcore/app/share"
	"gcore/dbmanager"
	"gcore/gcoordinator"
	"gcore/glog"
	"gcore/gmonitor/gmonitorService"
	"gcore/rpc"
	"os"
)

var logger = glog.NewLogger(`App`)

// GcoreApp - server instance struct
type GcoreApp struct {
	share       *share.AppShare
	coordinator gcoordinator.Coordinator
	monitorSrv  *gmonitorService.Service
	dbManager   *dbmanager.DBManager
	rpcManager  rpc.Manager
	signal      chan os.Signal

	startActs []func() error
	stopActs  []func() error
}

// GcoreAppConfig 启动配置信息
type GcoreAppConfig share.AppConfig

// GetCoordinator 返回协调模块对象
func (a *GcoreApp) GetCoordinator() gcoordinator.Coordinator {
	return a.coordinator
}

// GetRpc 返回rpc管理模块对象
func (a *GcoreApp) GetRpcManager() rpc.Manager {
	return a.rpcManager
}

// GetDBManager 返回数据库连接对象
func (a *GcoreApp) GetDBManager() *dbmanager.DBManager {
	return a.dbManager
}

// 信号控制接口
func (a *GcoreApp) Control(signal os.Signal) {
	app.signal <- signal
}

// 设置开机动作
func (a *GcoreApp) OnStart(starts ...func() error) {
	a.startActs = starts
}

// 设置关机动作
func (a *GcoreApp) OnStop(stops ...func() error) {
	a.stopActs = stops
}

// Boot 运行服务器消息监听循环，等待外部进程消息
func (a *GcoreApp) Boot(arg ...func()) {
	// 执行开机操作
	a.start()

	// 等待系统信号
	a.loop()

	// 执行关机操作
	a.stop()

	return
}

// 启动流程
func (a *GcoreApp) start() {
	logger.Infof("应用开始启动")

	conf := a.share.GetConfig()

	// 初始化分布式协调模块
	if conf.Has(`app/coordinator`) {
		logger.Infof("启动协调器")
		coord := gcoordinator.NewCoordinator(a.share, "app/coordinator/")

		err := coord.Start()
		if err == nil {
			app.coordinator = coord
		} else {
			logger.Fatalf("启动协调器出错: %s", err)
		}
	}

	// 如果有配置监控模块，则进行初始化
	if conf.Has(`app/monitorService`) {
		logger.Infof("启动指标监测服务")

		app.monitorSrv = gmonitorService.NewService(a.share, `app/monitorService/`)
		app.monitorSrv.Start()
	}

	// 初始化数据库连接
	if conf.Has(`app/dbManager`) {
		logger.Infof("启动数据库管理器")
		app.dbManager = dbmanager.NewManager(a.share, `app/dbManager/`)
	}

	// 初始化rpc
	if conf.Has(`app/coordinator`) {
		logger.Infof("启动 rpc 管理器")
		manager := rpc.NewManager(a.share, app.coordinator)

		err := manager.Start()
		if err == nil {
			app.rpcManager = manager
		} else {
			logger.Fatalf("启动 rpc 管理器出错: %s", err)
		}
	}

	// 执行业务开机动作
	for _, startAct := range a.startActs {
		err := startAct()
		if err != nil {
			logger.Fatalf("启动出错：%s", err)
		}
	}

	logger.Infof(`应用启动完成[PID=%d]`, os.Getpid())
}

// 关闭流程
func (a *GcoreApp) stop() {
	logger.Infof("应用开始关闭")

	var err error

	// 执行业务关闭动作
	for _, stopAct := range a.stopActs {
		err = stopAct()
		if err != nil {
			logger.Errorf("关闭出错: %s", err)
		}
	}

	conf := a.share.GetConfig()
	if conf.Has(`app/monitorService`) {
		logger.Infof("关闭指标监测服务")

		app.monitorSrv.Stop()
	}

	if conf.Has(`app/coordinator`) {
		logger.Infof("关闭 rpc 管理器")

		err = a.rpcManager.Stop()
		if err != nil {
			logger.Errorf("关闭 rpc 管理器出错: %s", err)
		}

		err = a.coordinator.Stop()
		if err != nil {
			logger.Errorf("关闭协调器出错: %s", err)
		}
	}

	logger.Infof("应用关闭完成")
}
