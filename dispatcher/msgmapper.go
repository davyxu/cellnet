package dispatcher

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	id2nameMap  = make(map[int]string)
	name2idMap  = make(map[string]int)
	mapperGuard sync.RWMutex
)

func addMapper(name string, id int) {

	mapperGuard.Lock()

	id2nameMap[id] = name
	name2idMap[name] = id

	mapperGuard.Unlock()
}

// 通过名字取id
func GetNameByID(id int) string {

	mapperGuard.RLock()
	defer mapperGuard.RUnlock()

	if v, ok := id2nameMap[id]; ok {
		return v
	}

	return "(??)"
}

// 通过id取名字
func GetIDByName(name string) int {

	mapperGuard.RLock()
	defer mapperGuard.RUnlock()

	if v, ok := name2idMap[name]; ok {
		return v
	}

	return 0
}

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
