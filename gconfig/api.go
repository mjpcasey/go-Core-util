package gconfig

import "gcore/gconfig/internal/pdspConfig"

// 构建应用配置
//
// @param file JSON 文件路径
func NewConfig(file string) Config {
	return pdspConfig.New(file)
}
