package utils

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func ProcessHost(host string) (string, error) {
	// var str string
	var ip string
	var port string

	// 是否包含“：”分割的字符串，取“：”后面的字符串做端口
	if strings.Contains(host, ":") {
		// 分割字符串
		arr := strings.Split(host, ":")
		if len(arr) > 2 {
			return "", errors.New("process host error, host invalid")
		}

		// 标准的ip：端口的配置, 例如：“192.168.9.102：9000”
		// 或者“:8000”的配置
		if len(arr) == 2 {
			port = arr[1]
			ip = arr[0]
		}
	} else {
		// 没有”：“分割，直接作为端口
		port = host
	}

	// 检查ip是否有设置，没有的话，取本机ip
	if ip == "" {
		ip = GetIp()
	}

	_, err := strconv.Atoi(port)
	if port == "" || err != nil {
		return "", errors.New("process host error, port invalid")
	}

	return ip + ":" + port, nil
}

func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var lip198, lip10, lip string
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "192.168.") == 0 {
				lip198 = sip
				break
			} else if strings.Index(addr.String(), "10.") == 0 {
				lip198 = sip
			} else if lip == "" && sip != "127.0.0.1" {
				lip = sip
			}
		}
	}
	if lip198 != "" {
		return lip198
	}
	if lip10 != "" {
		return lip10
	}
	return lip
}