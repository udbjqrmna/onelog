package onelog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRegisterInitRef(t *testing.T) {
	var initRun = func(patterns map[string]interface{}, writers map[string]interface{}) {
		patterns["aaa"] = &JsonPattern{}
	}

	err := RegisterInitRef(initRun)

	if err != nil {
		t.Error("操作对象为空")
	}
}

func Test2(t *testing.T) {
	var path = "./main/log.json"
	var config = make(map[string]interface{})

	f, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(NotFoundFile(path))
	}

	if err := json.Unmarshal(f, &config); err != nil {
		t.Error(err)
	}

	//生成对象以后的处理方式。处理缺省值
	//config.SetDefault()
	b, err := json.Marshal(config)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Println(string(b))

	fmt.Printf("%T,%s", config["WriterPara"].(map[string]interface{})["MaxCapacity"], config["WriterPara"].(map[string]interface{})["MaxCapacity"])
}

func TestNewLogFromConfig(t *testing.T) {
	var val float64 = 2
	var defLogLevel = strings.ToLower(Level(val).String())

	fmt.Println(defLogLevel)
	if Level(val) > Disable {
		fmt.Println("out ")
	} else {
		fmt.Println(Level(val))
	}

}

func TestConfig(t *testing.T) {
	LevelName = "l"

	if err := NewLogFromConfig("./main/log.json"); err != nil {
		t.Error(err)
		return
	}

	println(GetLog("one"), GetLog("two"), GetLog("three"))

	logOne := GetLog("one")

	logOne.Close()

	logTwo := GetLog("two")
	logTwo.Info().Msg("这是膛有一为人")
}

func TestConfigMultiple(t *testing.T) {
	LevelName = "l"

	if err := NewLogFromConfig("./main/log.json"); err != nil {
		fmt.Println(err)
	}

	logthree := GetLog("three")

	for i := 0; i < 500000; i++ {
		logthree.Info().Int("k",i).Msg("abcdefig")
	}

	logthree.Close()
}

func TestWriteFile2(t *testing.T) {
	f, err := createLogWriteFile("/Users/yimin/go/src/github.com/udbjqrmna/onelog/logs23/multiple.log")

	if err != nil {
		fmt.Println(err)
	}

	f.Write([]byte("00000000000"))

}
