package onelog

import (
	"time"
)

var logs = make(map[string]*Logger)

type Level uint8

const (
	TraceLevel Level = iota
	// DebugLevel defines debug log level.
	DebugLevel
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel

	Disable
)

var (
	MessageName     = "msg"
	CoroutineIDName = "cid"
	LevelName       = "level"
	TimeName        = "time"
	TimeFormat      = time.RFC3339
	CallerName      = "caller"
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}
	return ""
}

type Logger struct {
	lws      []LevelWriter
	writer   Writer
	minLevel Level
	pattern  Pattern
}

//NewLogger 返回一个新的Logger
func New(writer Writer, level Level, pattern Pattern) *Logger {
	var l = level
	if level > PanicLevel {
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
	for i := TraceLevel; i <= PanicLevel; i++ {
		if l.minLevel > i {
			l.lws[i] = disableLevelWriter
		} else { //如果已经设置过的将保留
			switch l.lws[i].(type) {
			case nil, *DisableLevelWriter:
				l.lws[i] = newDefaultLevelWriter(l.writer, i, l.pattern)
			}
		}
	}
}

//AddStatic 给此Logger所有日志都增加一个静态
func (l *Logger) AddStatic(name, value string) *Logger {
	//循环调用
	for i := TraceLevel; i <= PanicLevel; i++ {
		if l.lws[i] != disableLevelWriter {
			l.lws[i] = l.lws[i].AddStatic(name, value)
		}
	}

	return l
}

//AddRuntime 给此Logger所有日志都增加一个运行时记录
func (l *Logger) AddRuntime(r RunTimeCompute) *Logger {
	//循环调用
	for i := TraceLevel; i <= PanicLevel; i++ {
		if l.lws[i] != disableLevelWriter {
			l.lws[i] = l.lws[i].AddRuntime(r)
		}
	}

	return l
}

func (l *Logger) Close() {
	l.writer.Close()
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

//WarnLevel 返回一个Warn等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Warn() LevelWriter {
	if l.minLevel > WarnLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[WarnLevel].clone()
}

//ErrorLevel 返回一个Error等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Error() LevelWriter {
	if l.minLevel > ErrorLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[ErrorLevel].clone()
}

//FatalLevel 返回一个Fatal等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Fatal() LevelWriter {
	if l.minLevel > FatalLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[FatalLevel].clone()
}

//PanicLevel 返回一个Panic等级的日志对象。如果整体日志等级高于，则返回nil
func (l *Logger) Panic() LevelWriter {
	if l.minLevel > PanicLevel {
		return &DisableLevelWriter{}
	}

	return l.lws[PanicLevel].clone()
}

//SetLevelWriter 设置指定等级的LevelWriter对象，如果参数给的是nil.则会替换成DisableLevelWriter对象。
func (l *Logger) SetLevelWriter(level Level, leverWriter LevelWriter) *Logger {
	if leverWriter == nil {
		l.lws[level] = &DisableLevelWriter{}
	} else {
		l.lws[level] = leverWriter
	}

	return l
}

//SetLevel 设置Log的记录等级
func (l *Logger) SetLevel(level Level) *Logger {
	l.minLevel = level
	l.refresh()

	return l
}

//SaveLogList 将一个日志对象存入日志列表当中
func SaveLogList(name string, log *Logger) {
	logs[name] = log
}

//GetLog 将已经存入日志列表当中的日志对象取出,如果未找到将返回
func GetLog(name string) *Logger {
	return logs[name]
}

//InfoMsg 直接以Info等级进行一个日志记录
func (l *Logger) InfoMsg(msg string) {
	l.Info().Msg(msg)
}

//DebugMsg 直接以Debug等级进行一个日志记录
func (l *Logger) DebugMsg(msg string) {
	l.Debug().Msg(msg)
}

//TraceMsg 直接以Trace等级进行一个日志记录
func (l *Logger) TraceMsg(msg string) {
	l.Trace().Msg(msg)
}

//WarnMsg 直接以Warn等级进行一个日志记录
func (l *Logger) WarnMsg(msg string) {
	l.Warn().Msg(msg)
}

//ErrorMsg 直接以Error等级进行一个日志记录
func (l *Logger) ErrorMsg(msg string) {
	l.Error().Msg(msg)
}

//FatalMsg 直接以Fatal等级进行一个日志记录
func (l *Logger) FatalMsg(msg string) {
	l.Fatal().Msg(msg)
}

//PanicMsg 直接以Panic等级进行一个日志记录
func (l *Logger) PanicMsg(msg string) {
	l.Panic().Msg(msg)
}
