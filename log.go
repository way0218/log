package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	// LevelNone .
	LevelNone = -1
	// LevelMin .
	LevelMin = 0
	// LevelFatal .
	LevelFatal = 0
	// LevelError .
	LevelError = 1
	// LevelWarn .
	LevelWarn = 2
	// LevelInfo .
	LevelInfo = 3
	// LevelDebug .
	LevelDebug = 4
	// LevelTrace .
	LevelTrace = 5
	// LevelMax .
	LevelMax = 5
)

const (
	levelFatalName   string = "FATAL"
	levelErrorName   string = "ERROR"
	levelWarnName    string = "WARN "
	levelInfoName    string = "INFO "
	levelDebugName   string = "DEBUG"
	levelTraceName   string = "TRACE"
	levelInvalidName string = ""
)

// Logger .
type Logger struct {
	filename   string
	level      int
	rotateSize int
	currSize   int
	totalSize  int
	buf        []byte
	out        io.WriteCloser
	mutex      *sync.Mutex
}

// NewLogger .
func NewLogger(filename string, level int, rotateSize int, threadSafe bool) *Logger {
	var logger Logger
	logger.filename = filename
	logger.level = level

	if threadSafe {
		logger.mutex = &sync.Mutex{}
	}

	if filename == "stdout" {
		logger.out = os.Stdout
	} else if filename == "stderr" {
		logger.out = os.Stderr
	} else {
		var err error
		logger.out, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return nil
		}

		fileInfo, err := os.Stat(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return nil
		}

		logger.currSize = int(fileInfo.Size())
		logger.rotateSize = rotateSize
	}

	return &logger
}

// Close .
func (l *Logger) Close() {
	l.out.Close()
}

func (l *Logger) rotate() {
	l.Close()
	newFilename := l.filename + "." + time.Now().Format("20060102-150405")

	if err := os.Rename(l.filename, newFilename); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	var err error
	l.out, err = os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	l.currSize = 0
}

func (l *Logger) logv(level int, logStr string) int {
	_, fileName, line, ok := runtime.Caller(2)
	if !ok {
		fileName = "???"
		line = 0
	}

	shortFileName := fileName
	for i := len(fileName) - 1; i > 0; i-- {
		if fileName[i] == '/' {
			shortFileName = fileName[i+1:]
			break
		}
	}

	if l.mutex != nil {
		l.mutex.Lock()
	}

	l.buf = l.buf[:0]
	formatHeader(&l.buf, level)
	l.buf = append(l.buf, logStr...)

	l.buf = append(l.buf, " - "...)
	l.buf = append(l.buf, shortFileName...)
	l.buf = append(l.buf, ':')
	itoa(&l.buf, line, 0)

	l.buf = append(l.buf, '\n')

	logSize := len(l.buf)
	l.currSize += logSize
	l.totalSize += logSize

	l.out.Write(l.buf)

	if l.rotateSize > 0 && l.currSize > l.rotateSize {
		l.rotate()
	}

	if l.mutex != nil {
		l.mutex.Unlock()
	}

	return logSize
}

// Trace .
func (l *Logger) Trace(format string, v ...interface{}) int {
	if l.level < LevelTrace {
		return 0
	}

	return l.logv(LevelTrace, fmt.Sprintf(format, v...))
}

// Debug .
func (l *Logger) Debug(format string, v ...interface{}) int {
	if l.level < LevelDebug {
		return 0
	}

	return l.logv(LevelDebug, fmt.Sprintf(format, v...))
}

// Info .
func (l *Logger) Info(format string, v ...interface{}) int {
	if l.level < LevelInfo {
		return 0
	}

	return l.logv(LevelInfo, fmt.Sprintf(format, v...))
}

// Warn .
func (l *Logger) Warn(format string, v ...interface{}) int {
	if l.level < LevelWarn {
		return 0
	}

	return l.logv(LevelWarn, fmt.Sprintf(format, v...))
}

// Error .
func (l *Logger) Error(format string, v ...interface{}) int {
	if l.level < LevelError {
		return 0
	}

	return l.logv(LevelError, fmt.Sprintf(format, v...))
}

// Fatal .
func (l *Logger) Fatal(format string, v ...interface{}) int {
	if l.level < LevelFatal {
		return 0
	}

	return l.logv(LevelFatal, fmt.Sprintf(format, v...))
}

// SetLevel .
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// copy from package log
// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

func formatHeader(buf *[]byte, level int) {
	t := time.Now()
	year, month, day := t.Date()
	itoa(buf, year, 4)
	//*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	//*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e3, 6)

	levelStr := levelName(level)

	*buf = append(*buf, ' ')
	*buf = append(*buf, levelStr...)
	*buf = append(*buf, ' ')
}

func levelName(level int) string {
	switch level {
	case LevelFatal:
		return levelFatalName
	case LevelError:
		return levelErrorName
	case LevelWarn:
		return levelWarnName
	case LevelInfo:
		return levelInfoName
	case LevelDebug:
		return levelDebugName
	case LevelTrace:
		return levelTraceName
	default:
		return levelInvalidName
	}
}
