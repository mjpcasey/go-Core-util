package share

import (
	"gcore/gconfig"
)

// AppConfig - cargo/app Server init config struct
type AppConfig struct {
	File string
}

// AppShare - App base struct
type AppShare struct {
	Conf gconfig.Config
}

// Init 初始化服务器App对象
func (ca *AppShare) Init(cfg AppConfig) {
	ca.Conf = gconfig.NewConfig(cfg.File)
}

// GetConfig 返回默认配置对象
func (ca *AppShare) GetConfig() gconfig.Config {
	return ca.Conf
}

// Reload 重新加载配置信息
func (ca *AppShare) Reload() {
	ca.Conf.Reload()
}

// On 绑定应用事件回调
func (ca *AppShare) On(event string, callback func(data interface{})) int {
	return 0
}

// Emit 触发应用事件
func (ca *AppShare) Emit(event string, data interface{}) {

}
