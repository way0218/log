package log

import "fmt"

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger("stdout", LevelTrace, 0, true)
}

// OpenDefaultLog .
func OpenDefaultLog(filename string, level int, rotateSize int, threadSafe bool) bool {
	defaultLogger = NewLogger(filename, level, rotateSize, threadSafe)
	if defaultLogger == nil {
		fmt.Println("open default log fail, redirect to stdout")
		defaultLogger = NewLogger("stdout", LevelTrace, 0, true)
		return false
	}

	return true
}

// Trace .
func Trace(format string, v ...interface{}) {
	if defaultLogger.level < LevelTrace {
		return
	}

	defaultLogger.logv(LevelTrace, fmt.Sprintf(format, v...))
}

// Debug .
func Debug(format string, v ...interface{}) {
	if defaultLogger.level < LevelDebug {
		return
	}

	defaultLogger.logv(LevelDebug, fmt.Sprintf(format, v...))
}

// Error .
func Error(format string, v ...interface{}) {
	if defaultLogger.level < LevelError {
		return
	}

	defaultLogger.logv(LevelError, fmt.Sprintf(format, v...))
}

// Info .
func Info(format string, v ...interface{}) {
	if defaultLogger.level < LevelInfo {
		return
	}

	defaultLogger.logv(LevelInfo, fmt.Sprintf(format, v...))
}

// Warn .
func Warn(format string, v ...interface{}) {
	if defaultLogger.level < LevelWarn {
		return
	}

	defaultLogger.logv(LevelWarn, fmt.Sprintf(format, v...))
}

// Fatal .
func Fatal(format string, v ...interface{}) {
	if defaultLogger.level < LevelFatal {
		return
	}

	defaultLogger.logv(LevelFatal, fmt.Sprintf(format, v...))
}

// SetLevel .
func SetLevel(level int) {
	defaultLogger.level = level
}
