package onelog

import "unicode/utf8"

//Pattern 记录的格式
type Pattern interface {
	//init 初始化的动作
	init(buffer []byte) []byte
	//AppKey 将一个名称记录至缓存中
	AppendKey(buffer []byte, key string) []byte
	//AppendInt64 将一个int64值记录缓存中 base 指明此值以什么格式显示
	AppendInt64(buffer []byte, value int64, base int) []byte
	//AppendUint64 将一个uint64值记录缓存中 base 指明此值以什么格式显示
	AppendUint64(buffer []byte, value uint64, base int) []byte
	//AppendUint64 将一个uint32值记录缓存中 base 指明此值以什么格式显示
	AppendUint32(buffer []byte, value uint32, base int) []byte
	//AppendFloat64 将一个float64值记录缓存中
	AppendFloat64(buffer []byte, value float64) []byte
	//AppValue 将一个值string记录缓存中
	AppendString(buffer []byte, value string) []byte
	//AppValue 将一个值string记录缓存中
	AppendValue(buffer []byte, value []byte) []byte
	//complete 在整个记录完成时调用的方法，这个方法用来最后调整整个记录
	Complete(buffer []byte) []byte
	addRuntimeValues(buffer []byte, r RunTimeCompute) []byte
}

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

func appendStringComplex(dst []byte, s []byte, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if noEscapeTable[b] {
			i++
			continue
		}

		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRune(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}

		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}

	if start < len(s) {
		dst = append(dst, s[start:]...)
	}

	return dst
}
