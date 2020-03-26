package dbmanager

import (
	"errors"
	"fmt"
	"gcore/gmonitor"
	"time"

	"github.com/gocql/gocql"

	"gcore/app/share"

	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
)

// DB 管理器
type DBManager struct {
	appShare    *share.AppShare
	mongoMap    map[string]*mongoWrapper
	aeMap       map[string]*AeroSpikeWrapper
	mysqlMap    map[string]*MySQLWrapper
	redisMap    map[string]RedisWrapper
	influxMap   map[string]*InfluxDBWrapper
	scylladbMap map[string]ScyllaDBWrapper
}

// 构建DB管理器
func NewManager(app *share.AppShare, confPath string) *DBManager {
	manager := &DBManager{
		appShare:    app,
		mongoMap:    make(map[string]*mongoWrapper),
		aeMap:       make(map[string]*AeroSpikeWrapper),
		mysqlMap:    make(map[string]*MySQLWrapper),
		redisMap:    make(map[string]RedisWrapper),
		influxMap:   make(map[string]*InfluxDBWrapper),
		scylladbMap: make(map[string]ScyllaDBWrapper),
	}
	config := manager.appShare.GetConfig()
	if config.Has(confPath + `mongo`) {
		manager.createMongo(manager.appShare, confPath+`mongo/`)
	}
	if config.Has(confPath + `aerospike`) {
		manager.createAs(manager.appShare, confPath+`aerospike/`)
	}
	if config.Has(confPath + `mysql`) {
		manager.createMySQL(manager.appShare, confPath+`mysql/`)
	}
	if config.Has(confPath + `redis`) {
		manager.createRedis(manager.appShare, confPath+`redis/`)
	}
	if config.Has(confPath + `influxdb`) {
		manager.createInfluxdb(manager.appShare, confPath+`influxdb/`)
	}
	if config.Has(confPath + `scylladb`) {
		manager.createScyllaDB(manager.appShare, confPath+`scylladb/`)
	}
	return manager
}

// 获取 MongoDB 客户端封装
func (m *DBManager) GetMongo(form string) (db *mongoWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/mongo/0/name") {
		form = config.Get("app/dbManager/mongo/0/name")
	}

	db = m.mongoMap[form]
	if db == nil {
		err = errors.New("can't find form from mongoMap, form name: " + form)
		logger.Errorf(err.Error())
		return nil, err
	}
	return
}

// 获取 AeroSpike 客户端封装
func (m *DBManager) GetAeroSpike(form string) (db *AeroSpikeWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/aerospike/0/name") {
		form = config.Get("app/dbManager/aerospike/0/name")
	}

	db = m.aeMap[form]
	if db == nil {
		logger.Errorf("找不到 Aerospike 节点[%s]", form)
		return nil, err
	}
	return
}

// 获取 MySQL 客户端封装
func (m *DBManager) GetMySQL(form string) (db *MySQLWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/mysql/0/name") {
		form = config.Get("app/dbManager/mysql/0/name")
	}

	db = m.mysqlMap[form]
	if db == nil {
		err = fmt.Errorf("找不到 mysql 节点[%s]", form)

		logger.Errorf(err.Error())

		return nil, err
	}
	return
}

// 获取 Redis 客户端封装
func (m *DBManager) GetRedis(form string) (db RedisWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/redis/0/name") {
		form = config.Get("app/dbManager/redis/0/name")
	}

	db = m.redisMap[form]
	if db == nil {
		logger.Errorf("找不到 Redis 节点[%s]", form)

		return nil, fmt.Errorf("找不到 Redis 节点[%s]", form)
	}

	return
}

// 获取 InfluxDB 客户端封装
func (m *DBManager) GetInfluxdb(form string) (db *InfluxDBWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/influxdb/0/name") {
		form = config.Get("app/dbManager/influxdb/0/name")
	}

	db = m.influxMap[form]
	if db == nil {
		logger.Errorf("找不到 influxdb 节点[%s]", form)

		return nil, fmt.Errorf("找不到 influxdb 节点[%s]", form)
	}
	return
}

// 获取 ScyllaDB 客户端封装
func (m *DBManager) GetScyllaDB(form string) (db ScyllaDBWrapper, err error) {
	// 默认第一个数据库设置
	config := m.appShare.GetConfig()
	if form == "" && config.Has("app/dbManager/scylladb/0/name") {
		form = config.Get("app/dbManager/scylladb/0/name")
	}

	db = m.scylladbMap[form]
	if db == nil {
		logger.Errorf("找不到 ScyllaDB 节点[%s]", form)

		return nil, errors.New(fmt.Sprintf("找不到 ScyllaDB 节点[%s]", form))
	}

	return
}

// 创建 redis 客户端封装
func (m *DBManager) createRedis(app *share.AppShare, configPrefix string) {
	config := app.GetConfig()
	for index := 0; index < config.GetInt(configPrefix+"length"); index++ {
		wrapper := &redisWrapper{}

		err := config.Scan(fmt.Sprintf("%s%d", configPrefix, index), &wrapper.config)
		if err != nil {
			logger.Fatalf("加载 Redis 配置出错: %s", err)
		}

		logger.Infof("创建 redis 数据库封装[%s]", wrapper.config.Name)

		// 启动测试连接
		opt := &redis.ClusterOptions{
			Addrs:        wrapper.config.Addrs,
			Password:     wrapper.config.Password,
			DialTimeout:  time.Second * time.Duration(wrapper.config.DialTimeoutSec),
			ReadTimeout:  time.Second * time.Duration(wrapper.config.ReadTimeoutSec),
			WriteTimeout: time.Second * time.Duration(wrapper.config.WriteTimeoutSec),
			MaxRetries:   wrapper.config.MaxRetries,
		}
		wrapper.client = redis.NewClusterClient(opt)

		if pong, err := wrapper.Ping(); pong == "PONG" {
			m.redisMap[wrapper.config.Name] = wrapper
			// 设置服务监控
			gmonitor.NewRedisClusterMonitor(wrapper.client, wrapper.config.Name)
		} else {
			logger.Fatalf("Redis 初始化 PING 检查出错: %s", err)
		}
	}

	return
}

// 创建 mongo 客户端封装
func (m *DBManager) createMongo(app *share.AppShare, confPath string) {
	config := app.GetConfig()

	for i := 0; i < config.GetInt(confPath+"length"); i++ {
		wrapper := &mongoWrapper{}

		err := config.Scan(fmt.Sprintf("%s%d", confPath, i), &wrapper.config)
		if err == nil {
			if wrapper.config.IdCounter == "" {
				wrapper.config.IdCounter = DefaultIdCounter
			}
		} else {
			logger.Fatalf("解析mongodb配置[index=%s]失败: %s", i, err)
		}

		logger.Infof("创建 mongo 数据库封装[%s]", wrapper.config.Name)
		// 初始化mongo数据库
		session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Username:       wrapper.config.Username,
			Password:       wrapper.config.Password,
			Addrs:          wrapper.config.Addrs,
			Database:       wrapper.config.Database,
			ReplicaSetName: wrapper.config.ReplicaSet,
			Mechanism:      wrapper.config.Mechanism,

			Timeout:      time.Second * time.Duration(wrapper.config.TimeoutSec),
			PoolTimeout:  time.Second * time.Duration(wrapper.config.PoolTimeoutSec),
			ReadTimeout:  time.Second * time.Duration(wrapper.config.ReadTimeoutSec),
			WriteTimeout: time.Second * time.Duration(wrapper.config.WriteTimeoutSec),

			FailFast:      true,  // 断网时快速报错而不是等到超时
			MaxIdleTimeMS: 120e3, // 连接池内的连接闲置时间限制
		})

		if err == nil {
			wrapper.session = session.Copy()

			m.mongoMap[wrapper.config.Name] = wrapper

			// 设置服务监控
			gmonitor.NewMongoMonitor(wrapper.session, wrapper.config.Name)
		} else {
			logger.Fatalf("mongodb 初始化[name=%s]出错: %s", wrapper.config.Name, err)
		}
	}
}

// 创建 scyllaDB 客户端封装
func (m *DBManager) createScyllaDB(app *share.AppShare, configPrefix string) {
	config := app.GetConfig()
	for index := 0; index < config.GetInt(configPrefix+"length"); index++ {
		wrapper := &scyllaDBWrapper{}

		err := config.Scan(fmt.Sprintf("%s%d", configPrefix, index), wrapper)
		if err != nil {
			logger.Fatalf("加载 ScyllaDB 配置出错: %s", err)
		}

		logger.Infof("创建 scylladb 数据库封装[name=%s, port=%d, addrs=%+v]", wrapper.Name, wrapper.Port, wrapper.Hosts)

		cluster := gocql.NewCluster(wrapper.Hosts...)
		cluster.Port = wrapper.Port
		cluster.NumConns = wrapper.NumConns
		cluster.ConnectTimeout = time.Duration(wrapper.ConnectTimeoutSec) * time.Second
		cluster.Timeout = time.Duration(wrapper.TimeoutSec) * time.Second
		cluster.ReconnectInterval = time.Duration(wrapper.ReconnectIntervalSec) * time.Second

		if wrapper.Username != "" && wrapper.Password != "" {
			cluster.Authenticator = gocql.PasswordAuthenticator{
				Username: wrapper.Username,
				Password: wrapper.Password,
			}
		}

		scylladb, err := cluster.CreateSession()
		if err != nil {
			logger.Fatalf("ScyllaDB 初始化出错: %s", err)
		}

		wrapper.client = scylladb
		m.scylladbMap[wrapper.Name] = wrapper
	}

	return
}
