package onelog

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type NotNil string

func (e NotNil) Error() string {
	return string(e) + "不能为空"
}

type NotUnderstand string

func (e NotUnderstand) Error() string {
	return "不能理解的值:" + string(e)
}

type NotFoundFile string

func (e NotFoundFile) Error() string {
	return "未找到指定的文件:" + string(e)
}

type MistakeType struct {
	expected  string
	practical string
}

func (e *MistakeType) Error() string {
	return "错误的格式。预期:" + e.expected + "　实际:" + e.practical
}

//NewLogFromConfig 新创建日志对象从指定的配置文件当中
//并转码进入给出的config 对象当中
//path 配置文件所在路径
func NewLogFromConfig(path string) error {
	var config = make(map[string]interface{})

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return NotFoundFile(path)
	}

	if err := json.Unmarshal(f, &config); err != nil {
		return err
	}

	if err = completionConfig(config); err != nil {
		return err
	}

	if err = loadLogs(config); err != nil {
		return err
	}

	return nil
}

//completionConfig 补全使用缺省值的配置
func completionConfig(config map[string]interface{}) error {
	//设定全局LogLevel的值
	var defLogLevel = "debug"
	if val, ok := config["LogLevel"]; ok {
		switch val.(type) {
		//使用string形式
		case string:
			defLogLevel = strings.ToLower(val.(string))
			//使用数字的方式
		case float64:
			v := Level(val.(float64))
			if v > Disable || v < TraceLevel {
				return &MistakeType{"0..7", strconv.Itoa(int(val.(float64)))}
			}
			defLogLevel = strings.ToLower(Level(val.(float64)).String())
		default:
			return NotUnderstand("LogLevel")
		}
	}

	//设定全局Pattern的值
	var defPattern = "jsonpattern"
	if val, ok := config["Pattern"]; ok {
		switch val.(type) {
		case string:
			defPattern = strings.ToLower(val.(string))
		default:
			return NotUnderstand("Pattern")
		}
	}

	//设定全局Writer的值
	var defWriter = "console"
	if val, ok := config["Writer"]; ok {
		switch val.(type) {
		case string:
			defWriter = strings.ToLower(val.(string))
		default:
			return NotUnderstand("Writer")
		}
	}

	if err := checkCorrect("default", defLogLevel, defPattern, defWriter); err != nil {
		return err
	}

	//开始循环Logs，分别对每一个记录做配置
	if logs, ok := config["Logs"]; ok {
		switch logs.(type) {
		case []interface{}:
			for i, record := range logs.([]interface{}) {
				switch record.(type) {
				case map[string]interface{}:
					var rec = logs.([]interface{})[i].(map[string]interface{})

					if _, ok = rec["Id"]; !ok {
						return NotNil("ID")
					}
					if _, ok = rec["LogLevel"]; !ok {
						rec["LogLevel"] = defLogLevel
					} else {
						rec["LogLevel"] = strings.ToLower(rec["LogLevel"].(string))
					}

					if _, ok = rec["Pattern"]; !ok {
						rec["Pattern"] = defPattern
					} else {
						rec["Pattern"] = strings.ToLower(rec["Pattern"].(string))
					}
					if _, ok = rec["Writer"]; !ok {
						rec["Writer"] = defWriter
					} else {
						rec["Writer"] = strings.ToLower(rec["Writer"].(string))
					}

					if _, ok = rec["WriterPara"]; !ok {
						return NotNil("ID:" + rec["Id"].(string) + "的参数WriterPara")
					}

					if err := checkCorrect(rec["Id"].(string), rec["LogLevel"].(string), rec["Pattern"].(string), rec["Writer"].(string)); err != nil {
						return err
					}
				default:
					return NotUnderstand("数组内值必须为json")
				}
			}
		default:
			return NotUnderstand("Logs")
		}
	}

	return nil
}

//checkCorrect 判断所给值的正确性
func checkCorrect(id, logLevel, pattern, writer string) error {
	if _, ok := refLevel[logLevel]; !ok {
		return &MistakeType{"id:" + id + ",trace..disable", logLevel}
	}
	if _, ok := refPattern[pattern]; !ok {
		return NotUnderstand("id:" + id + ",Pattern")
	}
	if _, ok := refWriter[writer]; !ok {
		return NotUnderstand("id:" + id + ",Writer")
	}

	return nil
}

//loadLogs 从一个整理好的config里面获取值，并初始化logs对象
func loadLogs(config map[string]interface{}) error {
	if configs, ok := config["Logs"]; ok {
		for _, record := range configs.([]interface{}) {
			var r = record.(map[string]interface{})

			var pattern = reflect.New(reflect.TypeOf(refPattern[r["Pattern"].(string)])).Interface().(Pattern)
			var writer = reflect.New(reflect.TypeOf(refWriter[r["Writer"].(string)])).Interface().(Writer)

			//设置对象的实际参数
			if err := writer.SetConfig(r["WriterPara"].(interface{})); err != nil {
				return err
			}

			SaveLogList(r["Id"].(string), New(writer, refLevel["LogLevel"], pattern))
		}
	}

	return nil
}

var refLevel = make(map[string]Level)
var refPattern = make(map[string]interface{})
var refWriter = make(map[string]interface{})

func init() {
	//初始化反射Level对象
	refLevel["trace"] = TraceLevel
	refLevel["debug"] = DebugLevel
	refLevel["info"] = InfoLevel
	refLevel["warn"] = WarnLevel
	refLevel["error"] = ErrorLevel
	refLevel["fatal"] = FatalLevel
	refLevel["panic"] = PanicLevel
	refLevel["disable"] = Disable

	//初始化反射的WritePattern对象
	refPattern["jsonpattern"] = JsonPattern{}
	refPattern["old"] = OldPattern{}

	//初始化反射的Writer对象
	refWriter["console"] = Stdout{}
	refWriter["file"] = FileWriter{}
	refWriter["multiple"] = MultipleWriter{}

}

//RegisterInitRef 注册自己相关的初始化对象，此方法接受一个func(map[string]Pattern, map[string]Writer)参数。
//在func(map[string]Pattern, map[string]Writer)　方法内增加对应配置文件的名称与对象
func RegisterInitRef(init func(map[string]interface{}, map[string]interface{})) error {
	if init == nil {
		return NotNil("初始化方法　initRef")
	}

	init(refPattern, refWriter)

	return nil
}
