package util

import (
	"errors"
	"strings"
)

// 读取=分割的配置文件
func ReadKVFile(filename string, callback func(k, v string) bool) (ret error) {
	readErr := ReadFileLines(filename, func(line string) bool {

		line = strings.TrimSpace(line)

		// 注释
		if strings.HasPrefix(line, "#") {
			return true
		}

		// 等号切分KV
		pairs := strings.Split(line, "=")

		switch len(pairs) {
		case 1:
			value := strings.TrimSpace(pairs[0])

			if value == "" {
				return true
			}

			return callback("", value)
		case 2:
			key := strings.TrimSpace(pairs[0])
			value := strings.TrimSpace(pairs[1])

			if key == "" {
				return true
			}

			return callback(key, value)
		default:
			ret = errors.New("Require '=' splite key and value")
			return false
		}

	})

	if readErr != nil {
		return readErr
	}

	return
}

type KVPair struct {
	Key   string
	Value string
}

// 将=分割的文件按值读回
func ReadKVFileValues(filename string) (ret []KVPair, err error) {

	err = ReadKVFile(filename, func(k, v string) bool {

		ret = append(ret, KVPair{k, v})
		return true
	})

	return
}
