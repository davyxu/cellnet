package util

import (
	"testing"
)

func TestStrHash(t *testing.T) {
	if v := StringHash("gamedef.EnterGameREQ"); v != 0x28c2f4fb {
		t.Errorf("expect 0x28c2f4fb, got %x", v)
	}

	if v := StringHash("gamedef.EnterGameACK"); v != 0x43980899 {
		t.Errorf("expect 0x43980899, got %x", v)
	}
}
