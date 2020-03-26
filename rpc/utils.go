package rpc

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

func processHost(hostPort string) (string, error) {
	ip, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", fmt.Errorf("无效地址: %s", err)
	}

	if ip == "" {
		ip = getIp()

		logger.Infof("网络IP解析结果: %s", ip)
	}

	return net.JoinHostPort(ip, port), nil
}

func getIp() string {
	cn, err := net.DialTimeout("udp", "8.8.8.8:80", time.Second)

	if err == nil {
		ip, _, err := net.SplitHostPort(cn.LocalAddr().String())

		if err == nil {
			return ip
		}
	}

	logger.Fatalf("解析网络IP地址失败: %s", err)

	return ""
}

// 计步器，用于数组元素的顺序取用
type indexer struct {
	count uint64
}

func (idx *indexer) getIndex(len int) int {
	return int(atomic.AddUint64(&idx.count, 1) % uint64(len))
}
