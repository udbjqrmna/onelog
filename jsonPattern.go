package onelog

import (
	"math"
	"strconv"
)

//JsonPattern JSON的记录格式
type JsonPattern struct {
}

func (json *JsonPattern) init(buffer []byte) []byte {
	return append(buffer, '{')
}

//AppendKey 增加一个key的方法，key必须是一个string格式
func (json *JsonPattern) AppendKey(buffer []byte, key string) []byte {
	b := append(buffer, '"')
	b = appendStringComplex(b, []byte(key), 0)
	return append(b, '"', ':')
}

//AppendValue 增加一个[]byte数组值的方法
func (json *JsonPattern) AppendValue(buffer []byte, value []byte) []byte {
	d := append(buffer, value...)
	return append(d, ',')
}

//AppendUint64 将一个uint64值记录缓存中
func (json *JsonPattern) AppendUint64(buffer []byte, value uint64, base int) []byte {
	buffer = strconv.AppendUint(buffer, uint64(value), base)

	return append(buffer, ',')
}

//AppendUint64 将一个uint64值记录缓存中
func (json *JsonPattern) AppendUint32(buffer []byte, value uint32, base int) []byte {
	buffer = strconv.AppendUint(buffer, uint64(value), base)

	return append(buffer, ',')
}

//AppendFloat64 将一个float64值记录缓存中
func (json *JsonPattern) AppendFloat64(buffer []byte, val float64) []byte {
	switch {
	case math.IsNaN(val):
		return append(buffer, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(buffer, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(buffer, `"-Inf"`...)
	}

	strconv.AppendFloat(buffer, val, 'f', -1, 64)
	return append(buffer, ',')
}

//AppendInt64 将一个int64的值插入至数据内
func (json *JsonPattern) AppendInt64(buffer []byte, value int64, base int) []byte {
	buffer = strconv.AppendInt(buffer, value, base)

	return append(buffer, ',')
}

//AppendString 将一个string的值插入至数据内
func (json *JsonPattern) AppendString(buffer []byte, value string) []byte {
	buffer = append(buffer, '"')
	buffer = appendStringComplex(buffer, []byte(value), 0)
	buffer = append(buffer, '"', ',')

	return buffer
}

func (json *JsonPattern) Complete(buffer []byte) []byte {
	buffer[len(buffer)-1] = '}'
	return append(buffer, '\n')
}

func (json *JsonPattern) addRuntimeValues(buffer []byte, r RunTimeCompute) []byte {
	buffer = json.AppendKey(buffer, r.GetName())
	buffer = json.AppendValue(buffer, r.Values())

	return buffer
}