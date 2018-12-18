package log

import (
	"github.com/udbjqrmna/onelog"
	"os"
)

var log *onelog.Logger

func init() {
	log = onelog.New(&onelog.Stdout{os.Stdout}, onelog.DebugLevel, &onelog.JsonPattern{})
}

//GetLog 返回当前的log对象
func GetLog() *onelog.Logger {
	return log
}

//Trace 返回一个默认的Trace等级的日志对象。如果整体日志等级高于，则返回nil
func Trace() onelog.LevelWriter {
	return log.Trace()
}

//Debug 返回一个默认的Debug等级的日志对象。如果整体日志等级高于，则返回nil
func Debug() onelog.LevelWriter {
	return log.Debug()
}

//Info 返回一个默认的Info等级的日志对象。如果整体日志等级高于，则返回nil
func Info() onelog.LevelWriter {
	return log.Info()
}

//Fatal 返回一个默认的Fatal等级的日志对象。如果整体日志等级高于，则返回nil
func Fatal() onelog.LevelWriter {
	return log.Fatal()
}

//Error 返回一个默认的Error等级的日志对象。如果整体日志等级高于，则返回nil
func Error() onelog.LevelWriter {
	return log.Error()
}

//Warn 返回一个默认的Warn等级的日志对象。如果整体日志等级高于，则返回nil
func Warn() onelog.LevelWriter {
	return log.Warn()
}

//Panic 返回一个默认的Panic等级的日志对象。如果整体日志等级高于，则返回nil
func Panic() onelog.LevelWriter {
	return log.Panic()
}
