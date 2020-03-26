package region

import (
	"net"
	"path"
	"runtime"
	"testing"
)

var ipStrs []string
var errIpStrs []string
var ipv6Strs []string

func initTest() {

	var defaultConfig Config

	_, curFile, _, _ := runtime.Caller(1)
	defaultConfig.Ipv4Path = path.Dir(curFile) + "/geo.v6.2.dat"
	defaultConfig.Ipv6Path = path.Dir(curFile) + "/geo.v6.3.dat"

	ipStrs = []string{
		"255.255.255.255",
		"59.41.242.167",
		"180.149.134.17",
		"121.10.141.90",
		"210.14.69.133",
		"42.62.32.0",
	}

	ipv6Strs = []string{
		"2001:0DB8:0000:0000:ABCD:0000:0000:1234",
		"2002:DB8::ABCD:0:0:1234",
		"2001:DB8:0:0:ABCD::1234",
		"2409:8201:3400:0:0:0:0:0",
		"2409:8201:34ff:ffff:ffff:ffff:ffff:ffff",
		"2409:8232:0500:0:0:0:0:0",
	}

	errIpStrs = []string{
		"42.62.32.256",
		"",
		"aabbccc",
	}
	InitGeoBuffer(defaultConfig)

}

func TestToIpValue(t *testing.T) {
	// IPv4
	initTest()

	ipValue, _ := Ip2IpValue("0.0.0.0")

	if ipValue.IsIpv4 != true {
		t.Fatalf("is not ipv4")
	}

	// IPV6
	ipv6Str := "1050:0000:0000:0000:0005:0600:300c:326b"
	ipValue1, _ := Ip2IpValue(ipv6Str)

	if ipValue1.IsIpv4 == true {
		t.Fatalf("is not ipv6")
	}

	// IPV6短模式
	ipv6Str2 := "1050:0:0:0:5:600:300c:326b"
	ipValue2, _ := Ip2IpValue(ipv6Str2)
	if ipValue2.IsIpv4 == true {
		t.Fatalf("is not ipv6")
	}

	if CompareIpValue(ipValue1, ipValue2) != EQ {
		t.Fatalf("not support short format")
	}

}

func TestToString(t *testing.T) {

	items := []struct {
		ipStr string
		want  string
	}{
		{
			"2001:0DB8:0000:0000:ABCD:0000:0000:1234",
			"2001:db8::abcd:0:0:1234",
		},
		{
			"2001:DB8:0:0:ABCD::1234",
			"2001:db8::abcd:0:0:1234",
		},
		{
			"2001:DB8::ABCD:0:0:1234",
			"2001:db8::abcd:0:0:1234",
		},
	}

	for _, item := range items {
		ipstr, _ := ToString(item.ipStr)

		if ipstr != item.want {
			t.Fatalf("error: original str: %s ipstr: %s, want: %s", item.ipStr, ipstr, item.want)
		}
	}
}

func TestRead(t *testing.T) {
	initTest()

	// ipv6的测试
	ip1 := "2409:8201:3400:0:0:0:0:0"
	ip2 := "2409:8201:34ff:ffff:ffff:ffff:ffff:ffff"

	ipValue1, _ := Ip2IpValue(ip1)
	nCountryIdV1, nRegionIdV1, nCityIdV1, nDistrict1 := ReadIP(ipValue1)

	ipValue2, _ := Ip2IpValue(ip2)
	nCountryIdV2, nRegionIdV2, nCityIdV2, nDistrict2 := ReadIP(ipValue2)

	if nCountryIdV1 != nCountryIdV2 || nRegionIdV1 != nRegionIdV2 || nCityIdV1 != nCityIdV2 || nDistrict1 != nDistrict2 {
		t.Fatalf("un equal %d=%d %d=%d %d=%d, %d=%d\n", nCountryIdV1, nCountryIdV2, nRegionIdV1,
			nRegionIdV2, nCityIdV1, nCityIdV2, nDistrict1, nDistrict2)
	}

	if nCountryIdV1 != 10761 || nRegionIdV2 != 10782 || nCityIdV2 != 24328 || nDistrict1 != 0 {
		t.Fatalf("un equal: %d=%d %d=%d %d=%d, %d=%d", nCountryIdV1, 10761,
			nRegionIdV2, 10782, nCityIdV2, 24328, nDistrict1, 0)
	}

}

func TestEncodeIpDecodeIp(t *testing.T) {
	initTest()

	for _, ipStr := range ipStrs {
		encodeIp, err := EncodeIp(ipStr)
		if err != nil {
			t.Fatalf("encode err: ipstr=%s msg=%s", ipStr, err.Error())
		}

		t.Logf("encodeIp: %s", encodeIp)
		decodeIp, err := DecodeIp(encodeIp)
		if err != nil {
			t.Fatalf("decode err: ipstr=%s msg=%s", ipStr, err.Error())
		}

		t.Logf("decodeIp: %s", decodeIp)
		if net.ParseIP(ipStr).String() != decodeIp {
			t.Fatalf("unequal: %s=%s", ipStr, decodeIp)
		}
	}

	for _, ipStr := range errIpStrs {
		_, err := EncodeIp(ipStr)

		if err == nil {
			t.Fatalf("encode no err: ipstr=%s msg=%s", ipStr, err.Error())
		}

	}

}
