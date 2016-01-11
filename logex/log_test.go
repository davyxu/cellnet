package logex

import (
	"testing"
)

func TestLevel(t *testing.T) {

	logex := New("test")

	logex.Debugf("%d %s %v", 1, "hello", t)

	logex.Errorln("hello")

}
