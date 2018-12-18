package main

import (
	"fmt"
	"strconv"
	"time"
)

//一个实际应用的测试，正式使用时此文件应该不使用
//func main() {
//	var pattern = reflect.New(reflect.TypeOf(JsonPattern{})).Interface().(WritePattern)
//	var writer = reflect.New(reflect.TypeOf(Stdout{})).Interface().(Writer)
//	writer.SetConfig(nil)
//	//fmt.Printf("%T,%T",pattern,writer)
//
//	log := New(writer, TraceLevel, pattern)
//
//	log.Info().Msg("adiwpojfwef")
//	log.Error().Msg("一个错误")
//
//	log.Close()
//}

func main() {
	buf := make([]byte, 200)[0:0]

	buf = append(buf, "abcdefg.log"...)
	buf = append(buf, time.Now().Format(".0102_")...)
	buf = append(buf, strconv.Itoa(20)...)
	buf = append(buf, ".gz"...)

	fmt.Println(len(string(buf)), string(buf))

}

