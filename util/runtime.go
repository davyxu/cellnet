package util

import (
	"fmt"
	"runtime"
)

// 将绝对路径按需要的节从右压缩
// 例如: c:/Develop/nucleus.git/server/src/core/util/stack.go中
// 当cStripPathSection=3
// 返回 core/util/stack.go
func StripFileName(filename string, part int) string {

	slen := len(filename)
	slashCount := 0
	pos := 0
	for i := slen - 1; i >= 0; i-- {

		if filename[i] == '/' {
			slashCount++

			if slashCount >= part {
				pos = i
				break
			}
		}

	}

	// 如果确实找到有斜杠, 将斜杠的位置移除,表示这是一个相对路径
	if pos > 0 {
		pos++
	}

	return filename[pos:]
}

const cStripPathSection = 3

// 获取当前调用信息
func GetStackInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)

	if ok {
		return StripFileName(file, cStripPathSection), line
	}

	return "(unknown)", 0

}

func GetStackInfoString(skip int) string {
	file, line := GetStackInfo(skip + 1)
	return fmt.Sprintf("%s(%d)", file, line)
}
