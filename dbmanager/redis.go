package dbmanager

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// redis 封装接口
type RedisWrapper interface {
	GetNativeClient() (client *redis.ClusterClient, err error)
	// set 方法
	// 成功时返回 "OK"
	Set(key string, value string, expire time.Duration) (result string, err error)
	// get 方法
	// 当 key 不存在时，err有值（err == redis.Nil）
	// 注意判断
	Get(key string) (value string, err error)
	// ping 方法
	// 成功时返回 "PONG"
	Ping() (value string, err error)
	// exist 方法
	Exist(key string) (exist bool, err error)
}

// 接口实现
type redisWrapper struct {
	config struct {
		Name            string   `json:"name"`
		Addrs           []string `json:"addrs"`
		Password        string   `json:"password"`
		DialTimeoutSec  int      `json:"dialTimeoutSec"`
		ReadTimeoutSec  int      `json:"readTimeoutSec"`
		WriteTimeoutSec int      `json:"writeTimeoutSec"`
		MaxRetries      int      `json:"maxRetries"` // 注意不只网络错误造成重新请求算 retry，查询发生重定向也算 retry，也会被这个选项限制
	}

	client *redis.ClusterClient
}

// 获取 go-redis 原生客户端
func (r *redisWrapper) GetNativeClient() (client *redis.ClusterClient, err error) {
	if r.client == nil {
		err = errors.New("redis客户端没有初始化")
	}

	return r.client, err
}

// set 方法
// 成功时返回 "OK"
func (r *redisWrapper) Set(key string, value string, expire time.Duration) (result string, err error) {
	client, err := r.GetNativeClient()

	if err == nil {
		return client.Set(key, value, expire).Result()
	}

	return
}

// get 方法
// 当 key 不存在时，err有值（err == redis.Nil）
// 注意判断
func (r *redisWrapper) Get(key string) (value string, err error) {
	client, err := r.GetNativeClient()

	if err == nil {
		return client.Get(key).Result()
	}

	return
}

// ping 方法
// 成功时返回 "PONG"
func (r *redisWrapper) Ping() (value string, err error) {
	client, err := r.GetNativeClient()

	if err == nil {
		return client.Ping().Result()
	}

	return
}

// exist 方法
func (r *redisWrapper) Exist(key string) (exist bool, err error) {
	client, err := r.GetNativeClient()

	if err != nil {
		return
	}

	i, err := client.Exists(key).Result()

	return i == 1, err
}