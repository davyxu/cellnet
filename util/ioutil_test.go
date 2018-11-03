package util

import (
	"bytes"
	"testing"
)

type mywriter struct {
	bytes.Buffer
}

func (self *mywriter) Write(p []byte) (n int, err error) {

	if len(p) > 2 {
		n = 2
	} else {
		n = len(p)
	}

	self.Buffer.Write(p[0:n])

	return n, nil
}

func TestWriteFull(t *testing.T) {

	var m mywriter

	WriteFull(&m, []byte{1, 2, 3, 4, 5})

	t.Log(m.Bytes())
}

func TestCompressBytes(t *testing.T) {

	data := []byte("hello")

	data, err := CompressBytes(data)
	if err != nil {
		t.FailNow()
	}

	t.Log(DecompressBytes(data))
}

func TestStackToString(t *testing.T) {
	t.Log(StackToString(5))
}
