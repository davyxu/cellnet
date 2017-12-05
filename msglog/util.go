package msglog

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

func MessageName(msg interface{}) string {

	meta := cellnet.MessageMetaByType(reflect.TypeOf(msg).Elem())
	if meta == nil {
		return ""
	}

	return meta.Name
}

func MessageID(msg interface{}) int {

	meta := cellnet.MessageMetaByType(reflect.TypeOf(msg).Elem())
	if meta == nil {
		return 0
	}

	return int(meta.ID)
}

func MessageToString(msg interface{}) string {

	if msg == nil {
		return ""
	}

	if stringer, ok := msg.(interface {
		String() string
	}); ok {
		return stringer.String()
	}

	return ""
}
