package onelog

import (
	"math"
	"strconv"
)

//OldPattern JSON的记录格式
type OldPattern struct {
}

func (old *OldPattern) init(buffer []byte) []byte {
	return buffer
}

//AppendKey 增加一个key的方法，key必须是一个string格式
func (old *OldPattern) AppendKey(buffer []byte, key string) []byte {
	b := append(buffer, '\t')
	b = appendStringComplex(b, []byte(key), 0)
	return append(b, ':')
}

//func append()

//AppendValue 增加一个[]byte数组值的方法
func (old *OldPattern) AppendValue(buffer []byte, value []byte) []byte {
	return append(buffer, value...)
}

//AppendUint64 将一个uint64值记录缓存中
func (old *OldPattern) AppendUint64(buffer []byte, value uint64, base int) []byte {
	return strconv.AppendUint(buffer, uint64(value), base)
}

//AppendUint64 将一个uint64值记录缓存中
func (old *OldPattern) AppendUint32(buffer []byte, value uint32, base int) []byte {
	return strconv.AppendUint(buffer, uint64(value), base)
}

//AppendFloat64 将一个float64值记录缓存中
func (old *OldPattern) AppendFloat64(buffer []byte, val float64) []byte {
	switch {
	case math.IsNaN(val):
		return append(buffer, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(buffer, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(buffer, `"-Inf"`...)
	}
	return strconv.AppendFloat(buffer, val, 'f', -1, 64)
}

//AppendInt64 将一个int64的值插入至数据内
func (old *OldPattern) AppendInt64(buffer []byte, value int64, base int) []byte {
	return strconv.AppendInt(buffer, value, base)
}

//AppendString 将一个string的值插入至数据内
func (old *OldPattern) AppendString(buffer []byte, value string) []byte {
	return appendStringComplex(buffer, []byte(value), 0)
}

func (old *OldPattern) Complete(buffer []byte) []byte {

	return append(buffer, '\n')
}

func (old *OldPattern) addRuntimeValues(buffer []byte, r RunTimeCompute) []byte {
	switch r.(type) {
	case *TimeValue:
		var val = r.Values()
		buffer = append(buffer, val...)
		copy(buffer[len(val):], buffer)
		copy(buffer, val)
	default:
		buffer = old.AppendKey(buffer, r.GetName())
		buffer = old.AppendValue(buffer, r.Values())
	}

	return buffer
}
