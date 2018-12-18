# onelog
    一个使用go的日志操作，当前底层使用zerolog做为底层。后续向可桥接其他底层log进发。



## 特性：

* 提供快速的简单的使用方式，根据配置文件可将日志写入文件当中。
* d sdsd
* 


## 安装

```bash
go get github.com/udbjqrmna/onelog/log
```
## 使用前：

需要修改sdk,当前操作方式：
1.在src/runtime/proc.go最后增加如下方法。

```go
func Goid() int64 {
  _g_ := getg()
  return _g_.goid
}
```

## 最简单的使用方式
```go
import github.com/udbjqrmna/onelog/log

func main(){
  log.Debug().Msg("Hello world.")
}
```
>log默认使用输出为　os.Stdout\
>log默认使用记录格式为　json\
>log默认使用记录等级为　debug


## 配置

在应用程序目录下建立 config 目录 ，并在目录内新增 log.json文件

## 使用


## 日志等级：

* `TraceLevel`　跟踪等级，此等级为最低级
* `DebugLevel`
* `InfoLevel`
* `WarnLevel`
* `ErrorLevel`
* `FatalLevel`
* `PanicLevel`
* `Disable`　此等级将禁止日志的记录


>部分写入将会使用缓存，如文件写入，为保证完整写入数据，请在程序关闭前调用log.Close()方法。