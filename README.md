# go项目核心组件

## 新组件要求
包命名注意避免给业务代码造成同名干扰；取名为常用单词时，建议加前缀字母g

## 目录说明
## [app](./app)
应用包，用于创建应用实例，加载和管理功能服务，并提供数据库和配置接口。

## [dbmanager](./dbmanager)
数据库实例管理器组件

## [gconfig](./gconfig)
配置管理器组件

## [gcoordinator](./gcoordinator)
协调器组件

## [ghttp](./ghttp)
http 处理封装

## [glog](./glog)
日志组件

## [gmonitor](./gmonitor)
应用依赖服务和性能监控组件

## [goss](./goss)
对象存储读取工具

## [rpc](./rpc)
rpc 组件

## [utils](./utils)
其它工具