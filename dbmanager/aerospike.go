package dbmanager

import (
	"gcore/app/share"
	"gcore/glog"
	"strconv"

	as "github.com/aerospike/aerospike-client-go"
)

var logger = glog.NewLogger("dbmanager")

// 连接池大小
const defaultQueueSize = 1000

// AeroSpike 封装
type AeroSpikeWrapper struct {
	Name   string
	Addr   string
	Port   int
	Client *as.Client
}

// 创建 AeroSpike 封装
func (m *DBManager) createAs(app *share.AppShare, configPrefix string) {
	logger.Infof("创建 aerospike 数据库封装")

	config := app.GetConfig()
	for index := 0; index < config.GetInt(configPrefix+"length"); index++ {
		prefix := configPrefix + strconv.Itoa(index)
		wrapper := &AeroSpikeWrapper{
			config.Get(prefix + "/name"),
			config.Get(prefix + "/addr"),
			config.GetInt(prefix + "/port"),
			nil}

		policy := as.NewClientPolicy()
		policy.ConnectionQueueSize = config.GetIntDef(prefix+"/queueSize", defaultQueueSize)
		client, err := as.NewClientWithPolicy(policy, wrapper.Addr, wrapper.Port)
		if err != nil {
			logger.Infof("Aerospike 启动错误: %s", err)
			continue
		}
		wrapper.Client = client
		m.aeMap[wrapper.Name] = wrapper
	}
}
