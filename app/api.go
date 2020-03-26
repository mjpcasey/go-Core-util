package app

import (
	"gcore/app/share"
	"gcore/gconfig"
	"os"
	"path/filepath"
)

var app *GcoreApp

// 创建应用
func CreateApp(cfg string) *GcoreApp {
	if app == nil {
		logger.Infof("创建应用")
		if cfg == "" {
			cfg, _ = filepath.Abs("./conf/dev.json")
		}

		app = &GcoreApp{
			share:  &share.AppShare{},
			signal: make(chan os.Signal, 1),
		}

		for _, env := range []string{"http_proxy", "https_proxy", "HTTP_PROXY", "HTTPS_PROXY"} {
			if proxy, exist := os.LookupEnv(env); exist {
				logger.Warnf("网络代理[%s=%s]", env, proxy)
			}
		}

		wd, _ := os.Getwd()
		logger.Infof("工作目录[%s]", wd)
		logger.Infof("加载配置[%s]", cfg)
		app.share.Init(share.AppConfig(GcoreAppConfig{File: cfg}))
	} else {
		logger.Fatalf(`应用已创建`)
	}

	return app
}

func GetAppShare() *share.AppShare {
	if app == nil || app.share == nil {
		panic(`Server not inited`)
	}

	return app.share
}

// 获取应用单例
func GetApp() *GcoreApp {
	if app == nil {
		panic(`Server not inited`)
	}

	return app
}

// 获取应用配置单例
func GetConfig() gconfig.Config {
	return app.share.GetConfig()
}

// 结束
func Exit() {
	if app == nil {
		panic("应用未创建")
	}

	app.Control(os.Interrupt)
}
