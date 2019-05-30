package util

import (
	"bufio"
	"io"
	"net"
	"os"
)

// 完整发送所有封包
func WriteFull(writer io.Writer, buf []byte) error {

	total := len(buf)

	for pos := 0; pos < total; {

		n, err := writer.Write(buf[pos:])

		if err != nil {
			return err
		}

		pos += n
	}

	return nil

}

// 读取文本文件的所有行
func ReadFileLines(filename string, callback func(line string) bool) error {

	f, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	reader := bufio.NewScanner(f)

	reader.Split(bufio.ScanLines)
	for reader.Scan() {

		if !callback(reader.Text()) {
			break
		}
	}

	return nil
}

// 检查文件是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 获取文件大小
func FileSize(name string) int64 {
	if info, err := os.Stat(name); err == nil {
		return info.Size()
	}

	return 0
}

// 判断网络错误
func IsEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}
