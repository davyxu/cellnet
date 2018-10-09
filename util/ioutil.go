package util

import (
	"bufio"
	"io"
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
