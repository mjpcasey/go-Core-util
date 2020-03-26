# 应用包

用于创建应用实例，加载和管理功能服务，并提供数据库和配置接口。


## 示例代码
```go
// 程序主入口
func main() {
	// 新建应用
	a := app.CreateApp(os.Getenv("CONF_FILE"))

	// 设置开机动作
	a.OnStart(
		configService.Start,    // 启动渠道配置加载
		bidlogService.Start,    // 启动竞价日志服务
		freqService.Start,      // 启动频次服务
		budgetService.Start,    // 启动预算服务
		crowdService.Start,     // 启动人群查询服务
		storeService.Start,     // 启动数据加载服务
		ossService.Start,       // 对象存储服务
		algorithmService.Start, // 启动算法数据
		//antiSpamService.Start,   // 启动实时反作弊服务
		pdbControlService.Start, // 启动流量控制服务
		mapperService.Start,     // 启动数据映射服务
		bidRpcService.Start,     // 启动竞价rpc服务
		queryService.Start,      // 启动查询服务
		detectionService.Start,  // 启动过滤信息服务
	)

	// 设置关机动作
	a.OnStop(
		configService.Stop,
		detectionService.Stop,
		queryService.Stop,
		bidRpcService.Stop,
		//antiSpamService.Stop,
		mapperService.Stop,
		pdbControlService.Stop,
		algorithmService.Stop,
		ossService.Stop,
		storeService.Stop,
		crowdService.Stop,
		budgetService.Stop,
		freqService.Stop,
		bidlogService.Stop,
	)

	a.Boot()
}
```

## 常用接口
```
app.CreateApp() // 创建应用
app.GetDBManager() // 获取DB管理器
app.GetConfig() // 获取应用配置
app.OnStart() // 设置应用开机操作
app.OnStop() // 设置应用关机操作
```
