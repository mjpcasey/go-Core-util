package maskyLogger

import (
	"fmt"
	golog "log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	DefaultKeepDay = 90 // 按天日志的默认保留天数
)

var (
	DefaultLowCallDepth      = 0
	DefaultAppenderCallDepth = 2
	DefaultLoggerCallDepth   = 3
)
var UseShortFile bool

//日志的输出接口
type Appender interface {
	Log(extendCallDepth int, level string, args ...interface{})
	Logln(extendCallDepth int, level string, args ...interface{})
	Logf(extendCallDepth int, level string, format string, args ...interface{})
	SetCallDepth(int)
}

type emptyAppender struct{}

func NewEmptyAppender() *emptyAppender                                                   { return &emptyAppender{} }
func (this *emptyAppender) Log(extendCallDepth int, level string, args ...interface{})   {}
func (this *emptyAppender) Logln(extendCallDepth int, level string, args ...interface{}) {}
func (this *emptyAppender) Logf(extendCallDepth int, level string, format string, args ...interface{}) {
}
func (this *emptyAppender) SetCallDepth(int) {}

type baseAppender struct {
	*golog.Logger
	Name      string
	CallDepth int
}

func (l *baseAppender) SetCallDepth(level int) {
	l.CallDepth = level
}

func (l *baseAppender) log(extendCallDepth int, level string, fmtFunc func(...interface{}) string, args ...interface{}) {
	v := make([]interface{}, 1, len(args)+1)
	v[0] = "[" + level + "] "
	v = append(v, args...)
	if l.CallDepth == 0 {
		l.CallDepth = DefaultAppenderCallDepth
	}
	// fmt.Println("=============================")
	// fmt.Println(DefaultLowCallDepth, l.CallDepth, extendCallDepth)
	// for i := 0; i < l.CallDepth+extendCallDepth+2; i++ {
	// 	_, name, line, _ := runtime.Caller(i)
	// 	fmt.Println(name, line)
	// }
	// fmt.Println("=============================")

	l.Output(DefaultLowCallDepth+l.CallDepth+extendCallDepth, fmtFunc(v...))
}

func (l *baseAppender) logf(extendCallDepth int, level string, fmtFunc func(string, ...interface{}) string, format string, args ...interface{}) {
	format = "[" + level + "] " + format
	if l.CallDepth == 0 {
		l.CallDepth = DefaultAppenderCallDepth
	}
	l.Output(DefaultLowCallDepth+l.CallDepth+extendCallDepth, fmtFunc(format, args...))
}

func (l *baseAppender) Log(extendCallDepth int, level string, args ...interface{}) {
	l.log(extendCallDepth, level, fmt.Sprint, args...)
}

func (l *baseAppender) Logln(extendCallDepth int, level string, args ...interface{}) {
	l.log(extendCallDepth, level, fmt.Sprintln, args...)
}

func (l *baseAppender) Logf(extendCallDepth int, level string, format string, args ...interface{}) {
	l.logf(extendCallDepth, level, fmt.Sprintf, format, args...)
}

type FileAppender struct {
	*baseAppender
	fileName string
}

func NewFileAppender(name string, fileName string) *FileAppender {
	l := &FileAppender{
		baseAppender: &baseAppender{
			Name:      name,
			CallDepth: 2,
		},
		fileName: fileName,
	}
	return l
}

func (l *FileAppender) Log(extendCallDepth int, level string, args ...interface{}) {
	l.lazyNewLogger()
	l.log(extendCallDepth, level, fmt.Sprint, args...)
}

func (l *FileAppender) Logln(extendCallDepth int, level string, args ...interface{}) {
	l.lazyNewLogger()
	l.log(extendCallDepth, level, fmt.Sprintln, args...)
}

func (l *FileAppender) Logf(extendCallDepth int, level string, format string, args ...interface{}) {
	l.lazyNewLogger()
	l.logf(extendCallDepth, level, fmt.Sprintf, format, args...)
}

func (l *FileAppender) lazyNewLogger() {
	if l.baseAppender.Logger == nil {
		os.MkdirAll(path.Dir(l.fileName), 0755)
		logFile, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			golog.Fatalln("log conf error:", err.Error())
			return
		}
		var tag int
		if UseShortFile {
			tag = golog.LstdFlags | golog.Lshortfile
		} else {
			tag = golog.LstdFlags | golog.Llongfile
		}

		defaultLogger := golog.New(logFile, "", tag)
		l.baseAppender.Logger = defaultLogger
	}
}

type DailyAppender struct {
	*FileAppender
	today     string //当天日期
	fileName  string //文件
	extension string //后缀名
	keepday   int    //历史日志保留天数，-1 表示不保存
	lock      sync.Mutex
}

func NewDailyAppenderEx(name, fileName, extension string, keepday int) *DailyAppender {
	var fname string = fileName
	if strings.HasSuffix(fileName, extension) {
		fname = fileName[:len(fileName)-4]
	}

	if keepday == 0 {
		keepday = DefaultKeepDay
	}
	var appender = &DailyAppender{
		fileName:  fname,
		extension: extension,
		keepday:   keepday,
	}
	appender.setLogger(name, time.Now().Format("20060102"))

	//更新文件
	go func() {
		for {
			now := time.Now()
			h, m, s := now.Clock()
			leave := 86400 - (h*60+m)*60 + s
			select {
			case <-time.After(time.Duration(leave) * time.Second):
				appender.setLogger(name, time.Now().Format("20060102"))
			}
			// 删除日志
			kd := appender.keepday
			if kd > 0 {
				for i := kd; i < kd+30; i++ {
					day := now.AddDate(0, 0, -1*i)
					oldfilename := appender.completeFilename(day.Format("20060102"))
					os.Remove(oldfilename)
				}
			}
		}
	}()

	return appender
}

func NewDailyAppender(name, fileName string, keepday int) *DailyAppender {
	return NewDailyAppenderEx(name, fileName, ".log", keepday)
}

func (self *DailyAppender) setLogger(name string, day string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.today == day {
		return nil
	}
	self.today = day

	self.FileAppender = NewFileAppender(name, self.completeFilename(day))
	return nil
}

func (self *DailyAppender) completeFilename(day string) string {
	return self.fileName + "_" + day + self.extension
}

type ConsoleAppender struct {
	*baseAppender
}

func NewConsoleAppender(name string) *ConsoleAppender {
	var tag int
	if UseShortFile {
		tag = golog.LstdFlags | golog.Lshortfile
	} else {
		tag = golog.LstdFlags | golog.Llongfile
	}

	a := &ConsoleAppender{
		baseAppender: &baseAppender{
			Logger:    golog.New(os.Stdout, "", tag),
			Name:      name,
			CallDepth: 2,
		},
	}
	return a
}

//把不同等级的信息打印到不同的Appender中
type LevelSeparationAppender struct {
	Name      string
	appenders map[string]Appender
}

func NewLevelSeparationAppender(name string) *LevelSeparationAppender {
	return &LevelSeparationAppender{
		Name:      name,
		appenders: make(map[string]Appender),
	}
}

func (this *LevelSeparationAppender) SetLevelAppender(level string, appender Appender) {
	this.appenders[level] = appender
}
func (this *LevelSeparationAppender) Log(extendCallDepth int, level string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Log(extendCallDepth+1, level, args...)
	}
}

func (this *LevelSeparationAppender) Logln(extendCallDepth int, level string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Logln(extendCallDepth+1, level, args...)
	}
}

func (this *LevelSeparationAppender) Logf(extendCallDepth int, level string, format string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Logf(extendCallDepth+1, level, format, args...)
	}
}

func (this *LevelSeparationAppender) SetCallDepth(level int) {
	for _, ap := range this.appenders {
		ap.SetCallDepth(level)
	}
}

func NewLevelSeparationDailyAppender(name string, fileName string, keepday int) *LevelSeparationAppender {
	l := NewLevelSeparationAppender(name)

	fname := fileName
	if strings.HasSuffix(fileName, ".log") {
		fname = fileName[:len(fileName)-4]
	}
	for _, level := range logLevelStringMap {
		if level == "ALL" || level == "NO" {
			continue
		}
		levelAppender := NewDailyAppenderEx(name+"_"+level, fname, "."+strings.ToLower(level), keepday)
		l.SetLevelAppender(level, levelAppender)
	}
	return l
}
