package log

import (
	"testing"
)

func TestLevel(t *testing.T) {

	Debugf("%d %s %v", 1, "hello", t)

	Errorln("hello")

}
