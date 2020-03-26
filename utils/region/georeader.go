//从二进制的数据中获取geo信息
/*ipv4 文件格式:
|4字节：第一个ip信息的起始地址|4字节：最后一个ip信息的起始地址|4字节：记录总数|
|4字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
|4字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
|4字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|


ipv6 文件格式:
|16字节：第一个ip信息的起始地址|4字节：最后一个ip信息的起始地址|4字节：记录总数|
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
*/

package region

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/influxdata/influxdb/pkg/mmap"
)

var iPv4Geo GeoReader
var iPv6Geo GeoReader

var ERROR_IP_FORMAT = errors.New("error ip format")

type Config struct {
	Ipv4Path string
	Ipv6Path string
}

type IpValue struct {
	High   uint64
	Low    uint64
	IsIpv4 bool
}

func InitGeoBuffer(config Config) (err error) {

	if len(config.Ipv4Path) > 0 {
		iPv4Geo, err = NewGeoReader(config.Ipv4Path, true)
		if err != nil {
			return
		}
	}

	if len(config.Ipv6Path) > 0 {
		iPv6Geo, err = NewGeoReader(config.Ipv6Path, false)
		if err != nil {
			return
		}
	}
	return
}

func ReleaseGeoBuffer() error {
	iPv4Geo.ReleaseGeoBuffer()
	iPv6Geo.ReleaseGeoBuffer()

	return nil
}

type GeoReader interface {
	Read(ipValue IpValue) (country, region, city, district int)
	Parse(ip string) (country, region, city, district int)
	ReleaseGeoBuffer() error
}

// 这个界限不需要理会ipv4 还是 ipv6
func ParseIP(ip string) (country, region, city, district int) {
	isIpv4, err := IsIpv4(ip)
	if err != nil {
		return
	}

	if isIpv4 {
		country, region, city, district = iPv4Geo.Parse(ip)
	} else {
		country, region, city, district = iPv6Geo.Parse(ip)
	}

	return
}

// 已经知道了是IPV4的IP调用这个会快一点, 少了一步判断是否是IPV4的操作
func ParseIPV4(ip string) (country, region, city, district int) {
	return iPv4Geo.Parse(ip)
}

func ParseIPV6(ip string) (country, region, city, district int) {
	return iPv6Geo.Parse(ip)
}

func ReadIP(ipValue IpValue) (country, region, city, district int) {
	if ipValue.IsIpv4 {
		country, region, city, district = iPv4Geo.Read(ipValue)
	} else {
		country, region, city, district = iPv6Geo.Read(ipValue)
	}

	return
}

func ReadIPV4(ipValue IpValue) (country, region, city, district int) {
	return iPv4Geo.Read(ipValue)
}

func ReadIPV6(ipValue IpValue) (country, region, city, district int) {
	return iPv6Geo.Read(ipValue)
}

func NewGeoReader(path string, isIpv4 bool) (gr GeoReader, err error) {
	data, err := mmap.Map(path, 0)
	if err != nil {
		return
	}

	if len(data) == 0 {
		err = fmt.Errorf("没有数据")
		return
	}

	if isIpv4 {
		gr = &GeoReaderIPv4{
			buffer:   data,
			totalLen: int32(binary.BigEndian.Uint32(data[8:])),
			start:    int32(binary.BigEndian.Uint32(data[:4])),
			end:      int32(binary.BigEndian.Uint32(data[4:8])),
		}
	} else {
		gr = &GeoReaderIPv6{
			buffer:   data,
			totalLen: int32(binary.BigEndian.Uint32(data[8:])),
			start:    int32(binary.BigEndian.Uint32(data[:4])),
			end:      int32(binary.BigEndian.Uint32(data[4:8])),
		}
	}
	return
}
