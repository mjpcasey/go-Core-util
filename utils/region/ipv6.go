package region

import (
	"encoding/binary"

	"github.com/influxdata/influxdb/pkg/mmap"
)

type GeoReaderIPv6 struct {
	buffer   []byte //二进制数据
	totalLen int32  //数据个数

	start int32
	end   int32
}

func (this *GeoReaderIPv6) ReleaseGeoBuffer() error {
	return mmap.Unmap(this.buffer)
}

func (this *GeoReaderIPv6) Parse(ip string) (country, region, city, district int) {
	if this == nil {
		return
	}

	ipValue, err := Ip2IpValue(ip)
	if err != nil {
		return
	}

	return this.Read(ipValue)
}

func (this *GeoReaderIPv6) Read(ipValue IpValue) (country, region, city, district int) {
	if this == nil {
		return
	}

	num := this.search(ipValue)
	if num >= this.totalLen || num < 0 {
		return
	}

	return this.getGeoInfo(this.getIndex(num))
}

// 找出第一个大于等于ipValue的位置
func (geo *GeoReaderIPv6) search(ipValue IpValue) int32 {
	var low, high int32

	low = 0
	high = geo.totalLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		tmpIpValue := geo.getIPByIndex(geo.getIndex(mid))
		if CompareIpValue(tmpIpValue, ipValue) == LT {
			low = mid + 1
		} else if CompareIpValue(tmpIpValue, ipValue) == GT {
			high = mid - 1
		} else {
			return mid
		}
	}

	return low
}

//num：第几个Ip数据
func (geo *GeoReaderIPv6) getIndex(num int32) int32 {
	return geo.start + num*IPV6_T_SIZE
}

func (geo *GeoReaderIPv6) getIPByIndex(index int32) IpValue {
	start := index + IP_OFFSET
	return ipv6toIpValue(geo.buffer[start : start+IPV6_SIZE])
}

/*
//从二进制的数据中获取geo信息
ipv6 文件格式:
|16字节：第一个ip信息的起始地址|4字节：最后一个ip信息的起始地址|4字节：记录总数|
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|
 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
|16字节：ip |2字节:countryId|2字节:regionId|4字节:cityId|4字节:districtId|

*/
func (geo *GeoReaderIPv6) getGeoInfo(index int32) (nCountryId, nRegionId, nCityId, nDistrictId int) {

	nCountryId = int(binary.BigEndian.Uint16(geo.buffer[index+IPV6_COUNTRY_OFFSET : index+IPV6_REGION_OFFSET]))
	nRegionId = int(binary.BigEndian.Uint16(geo.buffer[index+IPV6_REGION_OFFSET : index+IPV6_CITY_OFFSET]))
	nCityId = int(binary.BigEndian.Uint32(geo.buffer[index+IPV6_CITY_OFFSET : index+IPV6_DISTRICT_OFFSET]))
	nDistrictId = int(binary.BigEndian.Uint32(geo.buffer[index+IPV6_DISTRICT_OFFSET : index+IPV6_DISTRICT_OFFSET+District_SIZE]))

	return
}
