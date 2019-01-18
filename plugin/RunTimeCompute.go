package plugin

import (
	"runtime"
	"strconv"
)

var CoroutineIDName = "cid"

//使用修改源码方式来获得当前执行此日志时的协程ID
// 需要修改sdk,当前操作方式：
//在`src/runtime/proc.go`最后增加如下方法。
//
//func Goid() int64 {
//  _g_ := getg()
//  return _g_.goid
//}
type CoroutineIDBySrc struct {
}

func (cid *CoroutineIDBySrc) GetName() string {
	return CoroutineIDName
}

func (cid *CoroutineIDBySrc) Values() []byte {
	result := make([]byte, 0)

	result = strconv.AppendInt(result, runtime.Goid(), 10)
	return result
}
