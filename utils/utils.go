package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"unsafe"
)

// 元素是否在切片中
func InStrSlice(slice []string, search string) bool {
	for i := range slice {
		if slice[i] == search {
			return true
		}
	}

	return false
}

// 元素是否在切片中
func InIntSlice(slice []int, search int) bool {
	for i := range slice {
		if slice[i] == search {
			return true
		}
	}

	return false
}

func NotInIntSlice(slice []int, search int) bool {
	for i := range slice {
		if slice[i] == search {
			return false
		}
	}

	return true
}

func IPv4ToUint32(s string) (i uint32) {
	i = uint32(big.NewInt(0).SetBytes(net.ParseIP(s).To4()).Uint64())

	return
}

// 取交集方法，注意参数需要事先去重
func Intersection(a, b []int) (c []int) {
	c = make([]int, 0, Min(len(a), len(b)))

	if len(a)*len(b) < 1e3 {
		for ai := range a {
			for bi := range b {
				if a[ai] == b[bi] {
					c = append(c, a[ai])
				}
			}
		}
	} else {
		m := make(map[int]byte)
		for i := range a {
			m[a[i]] = 0
		}
		for i := range b {
			if _, exist := m[b[i]]; exist {
				c = append(c, b[i])
			}
		}
	}

	return
}

// 字符串拼接
func Join(strArr []string) string {
	copyByte := make([]byte, 0, 40)
	bl := 0
	for _, str := range strArr {
		bl += copy(copyByte[bl:], str)
	}
	return *(*string)(unsafe.Pointer(&copyByte))
}

// 对字符串进行md5哈希,
// 返回32位小写md5结果
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5ToByte(s []byte) []byte {
	h := md5.New()
	h.Write(s)
	return h.Sum(nil)
}

// 对字符串进行签名，
// 返回int64的签名结果。
func CreateSign(s string) int64 {
	i, _ := strconv.ParseInt(MD5(s)[8:23], 16, 64)
	return i
}

func CreateSignInt32(s string) int {
	i, _ := strconv.ParseInt(MD5(s)[8:15], 16, 64)
	return int(i)
}

// 通过两重循环过滤重复元素（时间换空间）
func RemoveRepByLoop(slice []int) (result []int) {
	for i := range slice {
		for j := range result {
			if slice[i] == result[j] {
				goto next // 重复了
			}
		}
		result = append(result, slice[i])
	next:
	}

	return
}

// 通过map主键唯一的特性过滤重复元素（空间换时间）
func RemoveRepByMap(slice []int) (result []int) {
	pLen, m := 0, map[int]bool{}
	for i := range slice {
		pLen, m[slice[i]] = len(m), true
		if len(m) > pLen { // 加入map后，map长度变化，则元素不重复
			result = append(result, slice[i])
		}
	}
	return
}

// 元素去重（效率第一）
func RemoveRep(slice []int) []int {
	if len(slice) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slice)
	} else {
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slice)
	}
}

func Min(nums ...int) (m int) {
	m = nums[0]
	for _, n := range nums {
		if n < m {
			m = n
		}
	}
	return
}
func Max(nums ...int) (m int) {
	m = nums[0]
	for _, n := range nums {
		if n > m {
			m = n
		}
	}
	return
}

// 类似rand.shuffle，但是只在数据空间内清洗所需要的部分数据
type UniqRand struct {
	pool []int
	pos  int
}

func NewUniqRand(length, limit int) *UniqRand {
	pool := make([]int, limit)
	for i := 0; i < limit; i++ {
		pool[i] = i
	}
	for i := 0; i < length; i++ {
		j := rand.Intn(limit)
		pool[i], pool[j] = pool[j], pool[i]
	}
	return &UniqRand{pool: pool}
}

func (u *UniqRand) Int() int {
	u.pos++
	return u.pool[u.pos%len(u.pool)]
}

func (u *UniqRand) Slice(n int) []int {
	if n > len(u.pool) {
		n = len(u.pool)
	}
	return u.pool[:n]
}
