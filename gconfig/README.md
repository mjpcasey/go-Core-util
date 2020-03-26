# 配置包

业务不应该直接使用此包，应使用 app.GetConfig() 获取单例，单例是以下接口的实现，接口定义了各类配置值读取方法
```go 
type Config interface {
	Has(key string) bool
	Get(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	GetDef(key, def string) string
	GetIntDef(key string, def int) int
	GetFloatDef(key string, def float64) float64
	GetBoolDef(key string, def bool) bool
	GetVersion() int64
	GetPath() string
	Reload() error
	Scan(key string, pointer interface{}) error
}
```