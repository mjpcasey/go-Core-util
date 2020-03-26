# 数据库实例管理器

常用接口

```go
manager = NewManager() // 获取数据库实例管理器
manger.GetMongo() // 获取 Mongo 数据库驱动实例
manger.GetAeroSpike() // 获取 AeroSpike 数据库驱动实例
manger.GetMySQL() // 获取 MySQL 数据库驱动实例
manger.GetRedis() // 获取 Redis 数据库驱动实例
manger.GetInfluxdb() // 获取 Influxdb 数据库驱动实例
manger.GetScyllaDB() // 获取 ScyllaDB 数据库驱动实例
```
