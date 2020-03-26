package dbmanager

import (
	"gcore/app/share"
	"github.com/influxdata/influxdb/client/v2"
	"strconv"
)

// influxDB 数据库封装
type InfluxDBWrapper struct {
	Name     string
	Username string
	Password string
	Addr     string
	Database string
	Client   client.Client
}

// 创建 influxDB 数据库封装
func (m *DBManager) createInfluxdb(app *share.AppShare, configPrefix string) {
	logger.Infof("创建 influxdb 数据库封装")

	config := app.GetConfig()
	for index := 0; index < config.GetInt(configPrefix+"length"); index++ {
		// 获取mysql配置
		wrapper := &InfluxDBWrapper{
			config.Get(configPrefix + strconv.Itoa(index) + "/name"),
			config.Get(configPrefix + strconv.Itoa(index) + "/username"),
			config.Get(configPrefix + strconv.Itoa(index) + "/password"),
			config.Get(configPrefix + strconv.Itoa(index) + "/addr"),
			config.Get(configPrefix + strconv.Itoa(index) + "/database"),
			nil}

		err := wrapper.createClient()
		if err != nil {
			logger.Infof("m influxdb init error:", err)
			continue
		}

		m.influxMap[wrapper.Name] = wrapper
	}
	return
}

// 创建客户端
func (w *InfluxDBWrapper) createClient() (err error) {
	w.Client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     w.Addr,
		Username: w.Username,
		Password: w.Password,
	})
	return err
}

// 写方法
func (w *InfluxDBWrapper) Write(points []*client.Point) (err error) {
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: w.Database,
	})

	bps.AddPoints(points)

	err = w.Client.Write(bps)

	return err
}
