package logger

import (
	"fmt"
	"strings"
)

// Уровни логирования
const (
	LogOff     int8 = 0
	LogError   int8 = 1
	LogWarning int8 = 2
	LogInfo    int8 = 3
	LogDebug   int8 = 4
	LogTrace   int8 = 5
)

var currentLogLevel int8 = LogInfo // Можно менять во время выполнения

// Улучшенная функция логирования
func log(level int8, category string, args ...interface{}) {
	if level <= currentLogLevel {
		prefix := getLogPrefix(level, category)
		fmt.Print(prefix)
		fmt.Println(args...)
	}
}

func getLogPrefix(level int8, category string) string {
	var levelStr, emoji string

	switch level {
	case LogError:
		levelStr, emoji = "ERROR", "❌"
	case LogWarning:
		levelStr, emoji = "WARN ", "⚠️"
	case LogInfo:
		levelStr, emoji = "INFO ", "ℹ️"
	case LogDebug:
		levelStr, emoji = "DEBUG", "🐛"
	case LogTrace:
		levelStr, emoji = "TRACE", "🔍"
	default:
		levelStr, emoji = "UNKNOWN", "❓"
	}

	return fmt.Sprintf("%s [%s] %s: ", emoji, levelStr, strings.ToUpper(category))
}

// Удобные хелперы
func ErrorLog(category string, args ...interface{}) { log(LogError, category, args...) }
func WarnLog(category string, args ...interface{})  { log(LogWarning, category, args...) }
func InfoLog(category string, args ...interface{})  { log(LogInfo, category, args...) }
func DebugLog(category string, args ...interface{}) { log(LogDebug, category, args...) }
func TraceLog(category string, args ...interface{}) { log(LogTrace, category, args...) }
