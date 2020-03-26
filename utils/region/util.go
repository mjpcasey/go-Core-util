package region

import (
	"encoding/binary"
	"encoding/hex"
	"net"
)

func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}

	ip = ip.To4()
	if ip == nil {
		return 0
	}

	return binary.BigEndian.Uint32(ip)
}

func Long2Ip(ipuint32 uint32) (res string) {
	ipb := make([]byte, 4)
	binary.BigEndian.PutUint32(ipb, ipuint32)
	ip := net.IPv4(ipb[0], ipb[1], ipb[2], ipb[3])
	res = ip.String()
	return
}

func ipv4toIpValue(ip []byte) IpValue {
	return IpValue{
		Low:    uint64(binary.BigEndian.Uint32(ip)),
		IsIpv4: true,
	}
}

func ipv6toIpValue(ip []byte) IpValue {
	// 大端模式：高位存储在低位
	return IpValue{
		Low:    binary.BigEndian.Uint64(ip[8:]),
		High:   binary.BigEndian.Uint64(ip[:8]),
		IsIpv4: false,
	}
}

func Ip2IpValue(ipstr string) (IpValue, error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return IpValue{}, ERROR_IP_FORMAT
	}

	if ipV4 := ip.To4(); ipV4 != nil {
		return ipv4toIpValue(ipV4), nil
	} else {
		return ipv6toIpValue(ip), nil
	}
}

func IsIpv4(ipStr string) (bool, error) {
	for i := 0; i < len(ipStr); i++ {
		c := ipStr[i]
		if c == '.' {
			return true, nil
		} else if c == ':' {
			return false, nil
		}
	}

	return false, ERROR_IP_FORMAT
}

// -1: ipValue1 < ipValue2
// 0: ipValue1 == ipValue2
// 1: ipValue1 > ipValue2
func CompareIpValue(ipValue1 IpValue, ipValue2 IpValue) int {

	if ipValue1.IsIpv4 {
		if ipValue1.Low > ipValue2.Low {
			return GT
		} else if ipValue1.Low == ipValue2.Low {
			return EQ
		} else {
			return LT
		}
	} else {
		if ipValue1.High > ipValue2.High {
			return GT
		} else if ipValue1.High < ipValue2.High {
			return LT
		} else {

			if ipValue1.Low > ipValue2.Low {
				return GT
			} else if ipValue1.Low == ipValue2.Low {
				return EQ
			} else {
				return LT
			}
		}
	}
}

// ipv6 decode 出来的是小写短格式
func DecodeIp(ipStr string) (ret string, err error) {

	var decodeIp []byte
	if decodeIp, err = hex.DecodeString(ipStr); err != nil {
		return
	}

	ip := net.IP(decodeIp)

	if ipV4 := ip.To4(); ipV4 != nil {
		ret = ipV4.String()
	} else if ipV6 := ip.To16(); ipV6 != nil {
		ret = ipV6.String()
	} else {
		err = ERROR_IP_FORMAT
	}

	return
}

func EncodeIp(ipstr string) (ret string, err error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return "", ERROR_IP_FORMAT
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		return hex.EncodeToString(ipv4), nil
	} else {
		return hex.EncodeToString(ip), nil
	}
}

// 主要是ipv6的地址 --> 短格式， RFC 5952（IPv6地址文本表示建议书）建议用小写字母表示IPv6地址。
// 2001:0DB8:0000:0000:ABCD:0000:0000:1234 --> 2001:db8::abcd:0:0:1234
// 2001:DB8:0:0:ABCD::1234 --> 2001:db8::abcd:0:0:1234
func ToString(ipstr string) (ret string, err error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return "", ERROR_IP_FORMAT
	}

	return ip.String(), nil
}
