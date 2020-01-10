package util

import (
	"testing"
)

func TestDetectPort(t *testing.T) {
	DetectPort("100~200/path", func(a *Address, port int) (interface{}, error) {
		if port != 100 {
			t.FailNow()
		}

		return nil, nil
	})

	DetectPort("scheme://host:100~200/path", func(a *Address, port int) (interface{}, error) {
		if port != 100 {
			t.FailNow()
		}

		return nil, nil
	})

}
