package maskyLogger

const (
	NOSET = iota
	ALL
	DEBUG
	LOG
	NOTICE
	WARN
	ERROR
	PANIC
	NO
	LENGTH
)

var LogLevelMap = map[string]int{
	"NO":     NO,
	"FATAL":  NO,
	"DEBUG":  DEBUG,
	"WARN":   WARN,
	"NOTICE": NOTICE,
	"LOG":    LOG,
	"INFO":   LOG,
	"ERROR":  ERROR,
	"PANIC":  PANIC,
	"ALL":    ALL,
}

var logLevelStringMap = map[int]string{
	NO:     "NO",
	DEBUG:  "DEBUG",
	WARN:   "WARN",
	NOTICE: "NOTICE",
	LOG:    "LOG",
	ERROR:  "ERROR",
	PANIC:  "PANIC",
	ALL:    "ALL",
}
