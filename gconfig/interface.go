package gconfig

// 配置接口定义
type Config interface {
	Has(key string) bool                         // 路径是否存在
	Get(key string) string                       // 读取字符串
	GetInt(key string) int                       // 读取整数
	GetFloat(key string) float64                 // 读取浮点数
	GetBool(key string) bool                     // 读取布尔值
	GetDef(key, def string) string               // 读取字符串
	GetIntDef(key string, def int) int           // 读取整数
	GetFloatDef(key string, def float64) float64 // 读取浮点数
	GetBoolDef(key string, def bool) bool        // 读取布尔值
	GetVersion() int64                           // 获取配置版本
	GetPath() string                             // 获取配置路径
	Reload() error                               // 配置重载
	Scan(key string, pointer interface{}) error  // 解析配置至结构体
}
