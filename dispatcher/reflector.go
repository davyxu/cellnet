package dispatcher

import (
	"fmt"
	"github.com/davyxu/cellnet"
)

type mapperReflector struct {
	oldr cellnet.ContentReflector
}

func (self *mapperReflector) Reflect(data interface{}) string {

	switch d := data.(type) {
	case *cellnet.Packet:
		return fmt.Sprintf("Packet %s(%d) size:%d", GetNameByID(int(d.MsgID)), d.MsgID, len(d.Data))
	}

	if self.oldr != nil {
		return self.oldr.Reflect(data)
	}

	return ""
}

func init() {

	cellnet.SetContentReflector(&mapperReflector{oldr: cellnet.GetContentReflector()})

}
