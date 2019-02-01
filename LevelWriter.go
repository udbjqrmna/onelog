package onelog

import 	"fmt"

type LevelWriter interface {
	Int(key string, value int) LevelWriter
	Hex(key string, value int) LevelWriter
	Int64(key string, value int64) LevelWriter
	Uint64(key string, value uint64) LevelWriter
	Uint(key string, value uint) LevelWriter
	String(key, value string) LevelWriter
	Float32(key string, value float32) LevelWriter
	Float64(key string, value float64) LevelWriter
	Bool(key string, b bool) LevelWriter
	Bytes(key string, bytes []byte) LevelWriter
	//Msg 进行一次日志的消息写入，必须调用此方法或msgf()方法才能正常写入日志内
	Msg(message string)
	//Msg 进行一次日志的消息写入，参数可参考fmt.Sprintf()方法。
	Msgf(message string, p ...interface{})

	clone() LevelWriter
	//AddRuntime 增加一个计算值，这个值在每一个日志等级独立，每一次记录日志都将重新计算并记录它
	AddRuntime(r RunTimeCompute) LevelWriter
	//AddStatic 在日志里面增加一个常量值，这个常量值在每一个日志等级独立，每一次记录日志都将进行记录它
	AddStatic(name, value string) LevelWriter
}

var TRUE = []byte("true")
var FALSE = []byte("false")

func newDefaultLevelWriter(writer Writer, level Level, pattern Pattern) *DefaultLevelWriter {
	lw := &DefaultLevelWriter{
		buffer:  make([]byte, 256),
		Pattern: pattern,
		Writer:  writer,
	}

	lw.buffer = lw.Pattern.init(lw.buffer[:0])
	lw.AddStatic(LevelName, level.String())
	lw.AddRuntime(&TimeValue{})

	return lw
}

type DefaultLevelWriter struct {
	buffer          []byte
	Pattern         Pattern
	Writer          Writer
	runtimeComputes *RunTimeComputes
	origin          *DefaultLevelWriter
}

func (lw *DefaultLevelWriter) AddRuntime(r RunTimeCompute) LevelWriter {
	or := lw.origin
	if or == nil {
		or = lw
	}

	if r != nil {
		rs := &RunTimeComputes{
			r,
			or.runtimeComputes,
		}

		lw.runtimeComputes = rs
		or.runtimeComputes = rs
	}

	return lw
}

func (lw *DefaultLevelWriter) AddStatic(name, value string) LevelWriter {
	or := lw.origin
	if or == nil {
		or = lw
	}

	or.buffer = lw.Pattern.AppendKey(or.buffer, name)
	or.buffer = lw.Pattern.AppendString(or.buffer, value)

	return or.clone()
}

func (lw *DefaultLevelWriter) Hex(key string, value int) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendInt64(lw.buffer, int64(value), 16)

	return lw
}

func (lw *DefaultLevelWriter) Bytes(key string, bytes []byte) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendValue(lw.buffer, bytes)

	return lw
}

func (lw *DefaultLevelWriter) clone() LevelWriter {
	result := &DefaultLevelWriter{
		buffer:          make([]byte, cap(lw.buffer)),
		Pattern:         lw.Pattern,
		Writer:          lw.Writer,
		runtimeComputes: lw.runtimeComputes,
		origin:          lw,
	}

	copy(result.buffer, lw.buffer[:len(lw.buffer)])
	result.buffer = result.buffer[:len(lw.buffer)]
	return result
}

func (lw *DefaultLevelWriter) Int(key string, value int) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendInt64(lw.buffer, int64(value), 10)

	return lw
}
func (lw *DefaultLevelWriter) Int64(key string, value int64) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendInt64(lw.buffer, int64(value), 10)

	return lw
}
func (lw *DefaultLevelWriter) Uint64(key string, value uint64) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendUint64(lw.buffer, value, 10)

	return lw
}
func (lw *DefaultLevelWriter) Uint(key string, value uint) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendUint64(lw.buffer, uint64(value), 10)

	return lw
}
func (lw *DefaultLevelWriter) String(key, value string) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendString(lw.buffer, value)

	return lw
}
func (lw *DefaultLevelWriter) Float32(key string, value float32) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendFloat64(lw.buffer, float64(value))

	return lw
}
func (lw *DefaultLevelWriter) Float64(key string, value float64) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	lw.buffer = lw.Pattern.AppendFloat64(lw.buffer, value)

	return lw
}
func (lw *DefaultLevelWriter) Bool(key string, b bool) LevelWriter {
	lw.buffer = lw.Pattern.AppendKey(lw.buffer, key)
	if b {
		lw.buffer = lw.Pattern.AppendValue(lw.buffer, TRUE)
	} else {
		lw.buffer = lw.Pattern.AppendValue(lw.buffer, FALSE)
	}

	return lw
}

func (lw *DefaultLevelWriter) Msg(message string) {
	buf := lw.buffer
	pattern := lw.Pattern

	if lw.runtimeComputes != nil {
		run := lw.runtimeComputes
		for ; run != nil; run = run.next {
			buf = pattern.addRuntimeValues(buf, run.curr)
		}
	}

	buf = pattern.AppendKey(buf, MessageName)
	buf = pattern.AppendString(buf, message)
	buf = pattern.Complete(buf)

	_, _ = lw.Writer.Write(buf)
}

func (lw *DefaultLevelWriter) Msgf(message string, p ...interface{}) {
	buf := lw.buffer
	pattern := lw.Pattern

	if lw.runtimeComputes != nil {
		run := lw.runtimeComputes
		for ; run != nil; run = run.next {
			buf = pattern.addRuntimeValues(buf, run.curr)
		}
	}

	buf = pattern.AppendKey(buf, MessageName)
	buf = pattern.AppendString(buf, fmt.Sprintf(message, p))
	buf = pattern.Complete(buf)

	_, _ = lw.Writer.Write(buf)
}

type DisableLevelWriter struct {
}

var disableLevelWriter = &DisableLevelWriter{}

func (dlw *DisableLevelWriter) Hex(key string, value int) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Int(key string, value int) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Int64(key string, value int64) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Uint64(key string, value uint64) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Uint(key string, value uint) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) String(key, value string) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Float32(key string, value float32) LevelWriter {
	return dlw
}
func (dlw *DisableLevelWriter) Float64(key string, value float64) LevelWriter {
	return dlw
}

func (dlw *DisableLevelWriter) Bool(key string, b bool) LevelWriter {
	return dlw
}

//Msg 什么都不干的一个东西，直接返回。
func (dlw *DisableLevelWriter) Msg(message string) {
}

func (dlw *DisableLevelWriter) Msgf(message string, p ...interface{}) {
}

func (dlw *DisableLevelWriter) clone() LevelWriter {
	return dlw
}

func (dlw *DisableLevelWriter) Bytes(key string, bytes []byte) LevelWriter {
	return dlw
}

func (dlw *DisableLevelWriter) AddRuntime(r RunTimeCompute) LevelWriter {
	return dlw
}

func (dlw *DisableLevelWriter) AddStatic(name, value string) LevelWriter {
	return dlw
}
