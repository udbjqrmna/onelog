# onelog
一个go的日志操作，当前底层使用zerolog做为底层。后续向可桥接其他底层log进发。



## 特性：

* 一个快速的、容易上手的全功能日志
* 使用配置文件，无需复杂的配置
* 日志可打印控制台也可直接写入文件当中
* 日志格式可根据自己需要指定
* 日志的格式与writer可扩展

## 关于执行速度

我们认为生产环境的日志主要用来写入文件，而不是输出至控制台。因此本日志主要关注写入文件时的效率，并对此进行了优化。


## 安装

```bash
go get -u github.com/udbjqrmna/onelog/log
```

## 使用前：

需要修改sdk,当前操作方式：
1.在`src/runtime/proc.go`最后增加如下方法。

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
  log.Debug().Bool("boolVal", true).Int("intVal", 10).Msg("Hello world.")
}
```
>`log`默认使用输出为　`os.Stdout`\
>`log`默认使用记录格式为　json\
>`log`默认使用记录等级为　debug\
>`Int()`方法可在日志当中记录一个int值，相似的还有`Hex()`、`String()`、`Float64()`、`Bool()`\
>每次记录日志需要以Msg()做为最后的结束，未使用将不会被写入


## 新建日志对象
使用 `New(writer Writer, level Level, pattern WritePattern)` 方法得到一个新的日志对象，此对象为线程安全对象，可直接在线程当中使用
```go
import github.com/udbjqrmna/onelog

func main(){
  var log = onelog.New(&onelog.Stdout{os.Stderr}, onelog.InfoLevel, &onelog.JsonPattern{})

  log.Debug().Bool("boolVal", true).Int("intVal", 10).Msg("Hello world.")

  log.Close()
}
```
>在使用了`FileWriter`时，需要在程序结束的时候对`log`进行`Close()`操作，用以保证将缓存的内容写入至文件内

## 使用配置文件
使用 `NewLogFromConfig(path string)`　方法进行配置调用，配置好的`log`对象可使用 `GetLog(name string)` 方法获得
### 代码
```go
import github.com/udbjqrmna/onelog

func main(){
　if err := onelog.NewLogFromConfig("./config/log.json"); err != nil {
    fmt.Println(err)
　}

　log := onelog.GetLog("three")

　log.Info().Int("k",i).Msg("abcdefig")

　log.Close()
}
```

### log.json 文件
```json
{
  "LogLevel": "Trace",
  "Pattern": "JsonPattern",
  "Writer": "file",
  "WriterPara": {
    "LogsRoot": "./logs",
    "FileName": "log.log",
    "MaxCapacity": 1
  },
  "Logs": [
    {
      "Id": "one",
      "LogLevel": "Trace",
      "Pattern": "JsonPattern",
      "Writer": "file",
      "WriterPara": {
        "LogsRoot": "./logs.abc/ce",
        "FileName": "log2.log",
        "MaxCapacity": 1
      }
    },
    {
      "Id": "two",
      "LogLevel": "Info",
      "Pattern": "old",
      "Writer": "console",
      "WriterPara": {
        "Console": "Stderr"
      }
    },
    {
      "Id": "three",
      "LogLevel": "Info",
      "Pattern": "JsonPattern",
      "Writer": "multiple",
      "WriterPara": [
        {
          "Writer": "console",
          "WriterPara": {
            "Console": "Stderr"
          }
        },
        {
          "Writer": "file",
          "WriterPara": {
            "LogsRoot": "./logs23",
            "FileName": "multiple.log",
            "MaxCapacity": 5
          }
        }
      ]
    }
  ]
}
```
>如果使用了自定义的Pattern或Writer对象时，需要使用`RegisterInitRef` 方法进行注册，这样才能正确的获得\
>根结构当中指定的`LogLevel`、`Pattern`、`Writer`为默认值，在`Logs`节当中有未指定的时候使用根当中。

## 其他

### 日志等级：

* `TraceLevel`　跟踪等级，此等级为最低级
* `DebugLevel`
* `InfoLevel`
* `WarnLevel`
* `ErrorLevel`
* `FatalLevel`
* `PanicLevel`
* `Disable`　此等级将禁止日志的记录

### 日志格式
提供`JsonPattern`与`OldPattern`两种日志的格式，当然也可自己指定定义的日志格式。
在`New()`方法或配置文件的`Pattern`值当中指定使用的

### 写入对象
提供`Stdout`与`FileWriter`、`MultipleWriter`三种写入方式。当然也可自己指定定义的写入。
### 日志通用项
可为每一个日志的每一个日志等级实现独立的通用项设置，通用项设置好之后，每次日志将都自动将通用项带上

```go
import github.com/udbjqrmna/onelog

func main(){
  fw, _ := onelog.NewFileWriter("./logs/a.log", 5000000)
  mul := onelog.NewMultipleWriter(fw, &onelog.Stdout{os.Stdout})

  var log = onelog.New(mul, onelog.TraceLevel, &onelog.JsonPattern{})

  onelog.LevelName = "L"
  onelog.TimeFormat = ""

　log.Debug().AddRuntime(&CoroutineID{}).AddStatic("static", "staticValue")
　log.Error().AddRuntime(&CoroutineID{}).AddRuntime(&Caller{	CallerSkipFrameCount:1}).AddStatic("a1", "b2")

  log.Error().Int("INT", 10).Bool("B", true).Msg("This is a ErrorLevel")
  log.Debug().Int("INT", 10).Bool("B", true).Msg("This is a DebugLevel")

  log.Close()
}
```

#### 静态通用项
使用`AddStatic()`方法为每一个日志的日志等级进行独立的设置。每个设置将不影响其他等级的数据。

#### 运行时通用项
* `CoroutineID` 对象将获得当前执行此日志时的协程ID，可设置`onelog.CoroutineIDName`的值来改变它的项目名称
* `Caller` 对象获得当前调用者的文件与行数。可设置`onelog.CallerName`的值来改变它的项目名称，也可设置`CallerSkipFrameCount`的值来跳过几次函数的调用。
* 时间的运行时通用项是日志对象构建时自动加入的。可以使用`onelog.TimeFormat`方式进行改变.如直接设置成“”将使用UNIX时间的long值。
>设置名称的代码需要放至log日志实例或`NewLogFromConfig`方法之前进行。因此`最简单的使用方式`无法变更日志项名称

### 日志项名称自定义
每个日志项默认的名称可进行使用，使用类似`onelog.LevelName = "L"`的方法进行修改。
>此设置代码需要放至log日志实例或`NewLogFromConfig`方法之前进行。因此`最简单的使用方式`无法变更日志项名称
