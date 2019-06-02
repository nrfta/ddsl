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
	LEVEL_DRY_RUN
	LEVEL_DEBUG
)

var (
	logLevel = LEVEL_DRY_RUN
	levelMap = map[LogLevel]string{
		LEVEL_ERROR: "ERROR",
		LEVEL_WARN: "WARN",
		LEVEL_INFO: "INFO",
		LEVEL_DRY_RUN: "DRY-RUN",
		LEVEL_DEBUG: "DEBUG",
	}
)

func SetLogLevel(level LogLevel) {
	assertLevelValid(level)
	logLevel = level
}

func SetLogLevelStr(level string) {
	switch strings.ToUpper(level) {
	case "ERROR":
		SetLogLevel(LEVEL_ERROR)
	case "WARN":
		SetLogLevel(LEVEL_WARN)
	case "INFO":
		SetLogLevel(LEVEL_INFO)
	case "DRY-RUN":
		SetLogLevel(LEVEL_DRY_RUN)
	case "DEBUG":
		SetLogLevel(LEVEL_DEBUG)
	default:
		panic("unknown log level " + level)
	}
}

func Error(s string, a ...interface{}) {
	Log(LEVEL_ERROR, s, a...)
}
func Warn(s string, a ...interface{}) {
	Log(LEVEL_WARN, s, a...)
}
func Info(s string, a ...interface{}) {
	Log(LEVEL_INFO, s, a...)
}
func DryRun(s string, a ...interface{}) {
	Log(LEVEL_DRY_RUN, s, a...)
}
func Debug(s string, a ...interface{}) {
	Log(LEVEL_DEBUG, s, a...)
}

func Log(level LogLevel, s string, a ...interface{}) {
	if !isValidLevel(level) {
		panic(fmt.Sprintf("invalid log level %d", level))
	}
	if level <= logLevel {
		msg := s
		if len(a) > 0 {
			msg = fmt.Sprintf(msg, a...)
		}
		fmt.Printf("[%s] %s\n", levelMap[level], msg)
	}
}

func isValidLevel(level LogLevel) bool {
	return level >= LEVEL_ERROR && level <= LEVEL_DEBUG
}

func assertLevelValid(level LogLevel) {
	if !isValidLevel(level) {
		panic(fmt.Sprintf("invalid log level %d", level))
	}
}

