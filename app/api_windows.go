// +build windows

package app

//  执行重载配置
func ReloadConfig() {
	logger.Infof("执行配置重载")

	if app == nil {
		panic("应用未创建")
	}

	panic("implement me")
}
