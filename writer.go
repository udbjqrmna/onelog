package onelog

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Writer interface {
	io.Writer
	Close()
	SetConfig(config interface{}) error
}

type MultipleWriter struct {
	Writer Writer
	Next   *MultipleWriter
}

type FileWriter struct {
	file        *os.File
	fileName    string
	maxCapacity int64
	saveIndex   int
	buffer      []byte
	len         int
	mutex       sync.Mutex
	date        string
}

type Stdout struct {
	Writer io.Writer
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	return s.Writer.Write(p)
}

func (*Stdout) Close() {
}

//SetConfig 设置相关参数
func (s *Stdout) SetConfig(config interface{}) error {
	if config != nil {
		switch config.(type) {
		case map[string]interface{}:
		default:
			return &MistakeType{"map[string]interface {} type", ""}
		}

		if val, ok := config.(map[string]interface{})["Console"]; ok {
			switch val.(type) {
			case string:
				switch strings.ToLower(val.(string)) {
				case "stderr":
					s.Writer = os.Stderr
					return nil
				case "stdout":
					s.Writer = os.Stdout
					return nil
				}
			default:
				return &MistakeType{"string type", ""}
			}
		}
	}

	return NotUnderstand("Console")
}

func NewMultipleWriter(writers ...Writer) *MultipleWriter {
	var m = &MultipleWriter{}
	for _, w := range writers {
		if m.Writer == nil {
			m.Writer = w
			continue
		}

		m = &MultipleWriter{
			Writer: w,
			Next:   m,
		}
	}
	return m
}

func (m *MultipleWriter) Write(p []byte) (n int, err error) {
	if m.Writer != nil {
		var curr = m
		for ; curr != nil; curr = curr.Next {
			if n, err = curr.Writer.Write(p); err != nil {
				return
			}
		}
	} else {
		return 0, NotNil("未找到对象")
	}

	return len(p), nil
}

func (m *MultipleWriter) Close() {
	var curr = m
	for ; curr != nil; curr = curr.Next {
		curr.Writer.Close()
	}
}

//SetConfig 设置相关参数
func (m *MultipleWriter) SetConfig(config interface{}) error {
	if config != nil {
		switch config.(type) {
		case []interface{}:
			for i, record := range config.([]interface{}) {
				switch record.(type) {
				case map[string]interface{}:
					var rec = config.([]interface{})[i].(map[string]interface{})

					writer, ok := rec["Writer"]
					if !ok {
						return NotNil("Writer")
					}
					writerPara, ok := rec["WriterPara"]
					if !ok {
						return NotNil("WriterPara")
					}

					w := reflect.New(reflect.TypeOf(refWriter[writer.(string)])).Interface().(Writer)
					w.SetConfig(writerPara)

					if m.Writer == nil {
						m.Writer = w
						m.Next = nil
					} else {
						mo := &MultipleWriter{
							Writer: w,
							Next:   m.Next,
						}
						m.Next = mo
					}
				default:
					return NotUnderstand("数组内值必须为json")
				}
			}
			return nil
		default:
			return &MistakeType{"[]Json type", ""}
		}
	}
	return NotUnderstand("MultipleWriter:WriterPara")
}

func NewFileWriter(fileName string, maxCapacity int64) (*FileWriter, error) {
	var file, err = createLogWriteFile(fileName)
	if err != nil {
		return nil, err
	}

	var b = make([]byte, 1*1024*1024)

	return &FileWriter{
		file:        file,
		fileName:    fileName,
		maxCapacity: maxCapacity,
		buffer:      b,
		len:         0,
		mutex:       sync.Mutex{},
	}, nil
}

//Close 当程序关闭时调用的操作
func (w *FileWriter) Close() {
	w.writeToDisk(true)
}

//SetConfig 设置相关参数
func (w *FileWriter) SetConfig(config interface{}) error {
	if config != nil {
		switch config.(type) {
		case map[string]interface{}:
		default:
			return &MistakeType{"map[string]interface {} type", ""}
		}
		var filePath = ""

		if val, ok := config.(map[string]interface{})["LogsRoot"]; ok {
			switch val.(type) {
			case string:
				filePath = strings.TrimRight(val.(string), "/")

				//如果未指定根路径，直接给当前路径
				if filePath == "" {
					filePath = "."
				}

				if err := checkDir(filePath); err != nil {
					return err
				}
			default:
				return &MistakeType{"string type", ""}
			}
		}

		if val, ok := config.(map[string]interface{})["FileName"]; ok {
			switch val.(type) {
			case string:
				if val == "" {
					return NotNil("FileName")
				}
				//将路径与文件名合并
				filePath = filePath + "/" + strings.TrimLeft(val.(string), "/")
			default:
				return &MistakeType{"string type", ""}
			}
		}

		if val, ok := config.(map[string]interface{})["MaxCapacity"]; ok {
			switch val.(type) {
			case float64:
				v := int(val.(float64))
				if v <= 0 {
					return &MistakeType{"大于0", strconv.Itoa(v)}
				}

				newF, err := NewFileWriter(filePath, int64(v)*1024*1024)
				if err != nil {
					return err
				}

				w.fileName = newF.fileName
				w.file = newF.file
				w.len = newF.len
				w.saveIndex = newF.saveIndex
				w.buffer = newF.buffer
				w.mutex = newF.mutex
				w.maxCapacity = newF.maxCapacity

				return nil
			default:
				return &MistakeType{"string type", ""}
			}
		}
	}

	return NotUnderstand("WriterPara")
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(p) > cap(w.buffer)-w.len {
		w.writeToDisk(false)
	}

	copy(w.buffer[w.len:], p)
	w.len += len(p)

	return len(p), nil
}

func (w *FileWriter) writeToDisk(isClose bool) {
	if this, e := os.Stat(w.file.Name()); e == nil {
		if this.Size() > w.maxCapacity {
			var date = time.Now().Format(".0102_")
			if date != w.date {
				w.date = date
				w.saveIndex = 0
			}

			buf := make([]byte, 200)[0:0]

			buf = append(buf, w.fileName...)
			buf = append(buf, w.date...)
			index := len(buf)

			for true {
				w.saveIndex++
				buf = strconv.AppendInt(buf[:index], int64(w.saveIndex), 10)

				if exists(string(buf) + ".gz") {
					continue
				}
				break
			}
			tempName := string(buf)

			_ = w.file.Close()
			_ = os.Rename(w.fileName, tempName)
			w.file, _ = createLogWriteFile(w.fileName)

			//判断是否最后的结束，如果是最后结束，在压缩结束后将文件关闭
			if isClose {
				gzipFile(tempName)
				_ = w.file.Close()
			} else {
				go gzipFile(tempName)
			}
		}
	}

	_, _ = w.file.Write(w.buffer[:w.len])
	w.len = 0
}

//gzipFile 压缩文件，并删除原有的文件
func gzipFile(fileName string) {
	_ = CompressFile(fileName, fileName+".gz")
	_ = os.Remove(fileName)
}

//exists 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func createLogWriteFile(name string) (*os.File, error) {
	f, err := os.OpenFile(name, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

//CompressFile 使用gzip压缩成gz
func CompressFile(fileName string, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("出错了，%v", err)
		return err
	}
	_, _ = io.Copy(gw, file)
	if err != nil {
		fmt.Printf("出错了，%v", err)
		return err
	}

	return nil
}

func checkDir(path string) error {
	if !exists(path) {
		var upperIndex = strings.LastIndexByte(path, '/')
		if upperIndex > 0 {
			checkDir(path[:upperIndex])
		}

		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
