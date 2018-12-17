package onelog

import (
	"compress/gzip"
	"fmt"
	"github.com/udbjqr/gopak/common"
	"io"
	"os"
	"strconv"
	"sync"
	"syscall"
)

type Writer interface {
	io.Writer
	Close()
}

type MultipleWriter struct {
	writer Writer
	next   *MultipleWriter
}

type FileWriter struct {
	file        *os.File
	fileName    string
	maxCapacity int64
	saveIndex   int
	buffer      []byte
	buflen      int
	mutex       sync.Mutex
}

type Stdout struct {
	io.Writer
}

func (Stdout) Close() {
}

func NewMultipleWriter(writers ...Writer) *MultipleWriter {
	var m = &MultipleWriter{}
	for _, w := range writers {
		if m.writer == nil {
			m.writer = w
			continue
		}

		m = &MultipleWriter{
			writer: w,
			next:   m,
		}
	}
	return m
}

func (m *MultipleWriter) Write(p []byte) (n int, err error) {
	if m.writer != nil {
		var curr = m
		for ; curr != nil; curr = curr.next {
			if n, err = curr.writer.Write(p); err != nil {
				return
			}
		}
	} else {
		return 0, common.NotNil("未找到对象")
	}

	return len(p), nil
}

func (m *MultipleWriter) Close() {
	var curr = m
	for ; curr != nil; curr = curr.next {
		curr.writer.Close()
	}
}

func NewFileWriter(fileName string, maxCapacity int64) (*FileWriter, error) {
	//TODO 判断一下路径是否存在，不存在需要去创建出来

	var file, err = createLogWriteFile(fileName)
	if err != nil {
		return nil, err
	}

	var b = make([]byte, 2*1024*1024)

	return &FileWriter{
		file:        file,
		fileName:    fileName,
		maxCapacity: maxCapacity,
		buffer:      b,
		buflen:      0,
	}, nil
}

//Close 当程序关闭时调用的操作
func (w *FileWriter) Close() {
	w.writeToDisk(true)
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(p) > cap(w.buffer)-w.buflen {
		w.writeToDisk(false)
	}

	copy(w.buffer[w.buflen:], p)
	w.buflen += len(p)

	return len(p), nil
}

func (w *FileWriter) writeToDisk(isClose bool) {
	if this, e := os.Stat(w.file.Name()); e == nil {
		if this.Size() > w.maxCapacity {
			w.saveIndex++
			_ = w.file.Close()
			tempName := w.fileName + strconv.Itoa(w.saveIndex)
			_ = os.Rename(w.fileName, tempName)
			w.file, _ = createLogWriteFile(w.fileName)

			//判断是否最后的结束，如果是最后结束，在压缩结束后将文件关闭
			if isClose {
				gzipFile(w, tempName)
				_ = w.file.Close()
			} else {
				go gzipFile(w, tempName)
			}
		}
	}

	_, _ = w.file.Write(w.buffer[:w.buflen])
	w.buflen = 0
}

//gzipFile 压缩文件，并删除原有的文件
func gzipFile(w *FileWriter, fileName string) {
	_ = CompressFile(fileName, w.fileName+strconv.Itoa(w.saveIndex)+".log.gz")
	_ = os.Remove(fileName)
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
