package util

// 字符串转为16位整形值
func StringHash(s string) (hash uint16) {

	for _, c := range s {

		ch := uint16(c)

		hash = hash + ((hash) << 5) + ch + (ch << 7)
	}

	return
}
