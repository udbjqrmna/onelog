package onelog

import (
	"os"
	"time"
)

//var logs = make(map[string]*Logger)

type Level uint8

const (
	TraceLevel Level = iota
	// DebugLevel defines debug log level.
	DebugLevel
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLeven
	// ErrorLevel defines error log level.
	ErrorLeven
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLeven

	Disable
)

var (
	MessageName          = "msg"
	CoroutineIDName      = "cid"
	LevelName            = "level"
	TimeName             = "time"
	TimeFormat           = time.RFC3339
	CallerName           = "caller"
	CallerSkipFrameCount = 0
)

var defaultLog *Logger

func init() {
	defaultLog = New(Stdout{os.Stdout}, TraceLevel, JsonPattern{})
}
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLeven:
		return "WARN"
	case ErrorLeven:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLeven:
		return "PANIC"
	}
	return ""
}

type Logger struct {
	lws      []LevelWriter
	writer   Writer
	minLevel Level
	pattern  WritePattern
	//context    []byte
}

//NewLogger 返回一个新的Logger
func New(writer Writer, level Level, pattern WritePattern) *Logger {
	var l = level
	if level > PanicLeven {
		l = Disable
	}

	var log = &Logger{
		lws:      make([]LevelWriter, 8),
		writer:   writer,
		minLevel: l,
		pattern:  pattern,
	}

	log.refresh()

	return log
}

func (l *Logger) refresh() {
	for i := TraceLevel; i <= PanicLeven; i++ {
		if l.minLevel > i {
			l.lws[i] = disableLevelWriter
		} else {
			//TODO 这里可以根据不同的等级给不同的默认值。
			lw := newDefaultLevelWriter(l.writer, i, l.pattern)
			l.lws[i] = lw
		}
	}
}

func (l *Logger) Close() {
	l.writer.Close()
}

//TODO 可修改默认值等等

//Trace 返回一个默认的Trace等级的日志对象。如果整体日志等级高于，则返回nil
func Trace() LevelWriter {
	return defaultLog.Trace()
}

//Debug 返回一个默认的Debug等级的日志对象。如果整体日志等级高于，则返回nil
func Debug() LevelWriter {
	return defaultLog.Debug()
}

//Info 返回一个默认的Info等级的日志对象。如果整体日志等级高于，则返回nil
func Info() LevelWriter {
	return defaultLog.Info()
}

//Fatal 返回一个默认的Fatal等级的日志对象。如果整体日志等级高于，则返回nil
func Fatal() LevelWriter {
	return defaultLog.Fatal()
}

//Error 返回一个默认的Error等级的日志对象。如果整体日志等级高于，则返回nil
func Error() LevelWriter {
	return defaultLog.Error()
}

//Warn 返回一个默认的Warn等级的日志对象。如果整体日志等级高于，则返回nil
func Warn() LevelWriter {
	return defaultLog.Warn()
}

//Panic 返回一个默认的Panic等级的日志对象。如果整体日志等级高于，则返回nil
func Panic() LevelWriter {
	return defaultLog.Panic()
}

//TraceLevel 返回一个Trace等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Trace() LevelWriter {
	if l.minLevel > TraceLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[TraceLevel].clone()
}

//DebugLevel 返回一个Debug等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Debug() LevelWriter {
	if l.minLevel > DebugLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[DebugLevel].clone()
}

//InfoLevel 返回一个INFO等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Info() LevelWriter {
	if l.minLevel > InfoLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[InfoLevel].clone()
}

//WarnLeven 返回一个Warn等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Warn() LevelWriter {
	if l.minLevel > WarnLeven {
		return &DisableLevelWriter{}
	}

	return l.lws[WarnLeven].clone()
}

//ErrorLeven 返回一个Error等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Error() LevelWriter {
	if l.minLevel > ErrorLeven {
		return &DisableLevelWriter{}
	}

	return l.lws[ErrorLeven].clone()
}

//FatalLevel 返回一个Fatal等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Fatal() LevelWriter {
	if l.minLevel > FatalLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[FatalLevel].clone()
}

//PanicLeven 返回一个Panic等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Panic() LevelWriter {
	if l.minLevel > PanicLeven {
		return &DisableLevelWriter{}
	}

	return l.lws[PanicLeven].clone()
}
