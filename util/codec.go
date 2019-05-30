package util

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
)

// 字符串转为16位整形哈希
func StringHash(s string) (hash uint16) {

	for _, c := range s {

		ch := uint16(c)

		hash = hash + ((hash) << 5) + ch + (ch << 7)
	}

	return
}

// 字节计算MD5
func BytesMD5(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// 字符串计算MD5
func StringMD5(str string) string {
	return BytesMD5([]byte(str))
}

// 压缩字节
func CompressBytes(data []byte) ([]byte, error) {

	var buf bytes.Buffer

	writer := zlib.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()

	return buf.Bytes(), nil
}

// 解压字节
func DecompressBytes(data []byte) ([]byte, error) {

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}
