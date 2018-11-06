package util

import (
	"errors"
	"testing"
)

func TestDetectPort(t *testing.T) {
	DetectPort("scheme://host:100~200/path", func(s string) (interface{}, error) {
		if s != "host:100" {
			t.FailNow()
		}

		return nil, errors.New("err")
	})
}
