package log

import (
	"fmt"
	"strings"
)

type LogLevel int

const (
	LEVEL_ERROR LogLevel = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG

)

var (
	logLevel = LEVEL_WARN
)

func SetLogLevel(level LogLevel) {
	if level < LEVEL_ERROR || level > LEVEL_DEBUG {
		panic("log level must be ERROR, WARN, INFO or DEBUG")
	}
}
func SetLogLevelStr(level string) {
	switch strings.ToUpper(level) {
	case "ERROR":
		SetLogLevel(LEVEL_ERROR)
	case "WARN":
		SetLogLevel(LEVEL_WARN)
	case "INFO":
		SetLogLevel(LEVEL_INFO)
	case "DEBUG":
		SetLogLevel(LEVEL_DEBUG)
	default:
		panic("unknown log level " + level)
	}
}

func Error(s string, a ...interface{}) {
	log("ERROR", s, a...)
}
func Warn(s string, a ...interface{}) {
	log("WARN", s, a...)
}
func Info(s string, a ...interface{}) {
	log("INFO", s, a...)
}
func Debug(s string, a ...interface{}) {
	log("DEBUG", s, a...)
}

func log(prefix, s string, a ...interface{}) {
	fmt.Printf("[" + prefix + "] " + s + "\n", a...)
}


