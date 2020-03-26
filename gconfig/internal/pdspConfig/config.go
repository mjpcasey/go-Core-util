package pdspConfig

import (
	"encoding/json"
	"fmt"
	"gcore/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var logger = glog.NewLogger("config")

type configNode struct {
	s string
	i int
	f float64
	b bool
	o interface{} //original value
}
type interfaceMap map[string]interface{}
type cacheMapType map[string]configNode

// config manager app json config object
type config struct {
	path      string
	cache     cacheMapType
	timestamp int64
}

// Reload - clear config cache and reload
func (c *config) Reload() error {
	return c.parseConfigFile()
}

// GetVersion 获取配置更新时间版本
func (c *config) GetVersion() int64 {
	return c.timestamp
}

// GetPath 获取配置新建路径
func (c *config) GetPath() string {
	return c.path
}

// return config data
func (c *config) getValue(key string) *configNode {
	val, found := c.cache[key]
	if found {
		return &val
	}
	return nil
}

// Has check the key if exisit in the config
func (c *config) Has(key string) bool {
	_, found := c.cache[key]
	return found
}

// Get one config value string
// config.Get("subConfig/name")
func (c *config) Get(key string) string {
	return c.GetDef(key, "")
}

// GetInt config value as interget
func (c *config) GetInt(key string) int {
	return c.GetIntDef(key, 0)
}

// GetFloat config value as interget
func (c *config) GetFloat(key string) float64 {
	return c.GetFloatDef(key, 0)
}

// GetBool config value as bool
func (c *config) GetBool(key string) bool {
	return c.GetBoolDef(key, false)
}

// GetDef - return config value string with default value
func (c *config) GetDef(key, def string) string {
	val := c.getValue(key)
	if val != nil {
		return val.s
	}
	return def
}

// GetIntDef - return config value int with default value
func (c *config) GetIntDef(key string, def int) int {
	val := c.getValue(key)
	if val != nil {
		return val.i
	}
	return def
}

// GetFloatDef - return config value int with default value
func (c *config) GetFloatDef(key string, def float64) float64 {
	val := c.getValue(key)
	if val != nil {
		return val.f
	}
	return def
}

func (c *config) GetBoolDef(key string, def bool) bool {
	val := c.getValue(key)
	if val != nil {
		return val.b
	}
	return def
}

func (c *config) Scan(key string, pointer interface{}) error {
	val := c.getValue(key)
	if val == nil {
		return fmt.Errorf("配置 %s 不存在", key)
	}

	data, err := json.Marshal(val.o)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, pointer)
}

func (c *config) LoadExtendConf(path string, pointer interface{}) error {
	panic("等待实现")
}

func New(file string) (c *config) {
	path, err := filepath.Abs(file)
	if err != nil {
		logger.Fatalf("加载配置错误: %s", err)
	}

	c = &config{path: path}

	err = c.parseConfigFile()
	if err != nil {
		logger.Fatalf("解析配置错误: %s", err)
	}

	return
}

// read and parse config file data
func (c *config) parseConfigFile() error {
	jsonBytes, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	// 替换环境变量
	jsonString := replaceVariable(string(jsonBytes))

	var jsonCache interfaceMap
	if err := json.Unmarshal([]byte(jsonString), &jsonCache); err != nil {
		return err
	}
	// process config data into array
	cache := make(cacheMapType)

	translateValue(&cache, "", jsonCache)

	c.cache = cache
	// 日志更新时间
	c.timestamp = time.Now().Unix()
	return nil
}

func translateValue(c *cacheMapType, prefix string, val interface{}) {
	v := reflect.ValueOf(val)
	node := configNode{}

	switch v.Kind() {
	case reflect.Map:
		if len(prefix) > 0 {
			node.o = val
			(*c)[prefix] = node
			prefix += "/"
		}
		for _, key := range v.MapKeys() {
			translateValue(c, prefix+key.String(), v.MapIndex(key).Interface())
		}
		return
	case reflect.Slice:
		if len(prefix) > 0 {
			node.o = val
			(*c)[prefix] = node
			prefix += "/"
		}
		(*c)[prefix+"length"] = configNode{i: v.Len()}
		for idx := 0; idx < v.Len(); idx++ {
			translateValue(c, prefix+strconv.Itoa(idx), v.Index(idx).Interface())
		}
		return
	case reflect.String:
		node.s = val.(string)
		node.i, _ = strconv.Atoi(node.s)
		node.f, _ = strconv.ParseFloat(node.s, 64)
		node.o = val
	case reflect.Float32:
		f32 := val.(float32)
		node.f = float64(f32)
		node.i = int(f32)
		node.s = strconv.FormatFloat(node.f, 'f', -1, 32)
		node.o = val
	case reflect.Float64:
		node.f = val.(float64)
		node.i = int(node.f)
		node.s = strconv.FormatFloat(node.f, 'f', -1, 64)
		node.o = val
	case reflect.Int:
		node.i = val.(int)
		node.f = float64(node.i)
		node.s = strconv.Itoa(node.i)
		node.o = val
	case reflect.Bool:
		node.b = val.(bool)
		node.o = val
	}
	(*c)[prefix] = node
}

// replace ${var_name} macro string
var regx = regexp.MustCompile(`\${[A-Za-z0-9\-_]+}`)

func replaceVariable(str string) string {
	return regx.ReplaceAllStringFunc(str, replaceCallback)
}
func replaceCallback(match string) string {
	if `ENV_` == match[2:6] {
		// if env, found := os.LookupEnv(match[6 : len(match)-1]); found {
		// 	return env
		// }
		match = os.Getenv(match[6 : len(match)-1])
	}
	return match
}
