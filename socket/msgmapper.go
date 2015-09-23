package socket

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

var (
	id2nameMap  = make(map[int]string)
	name2idMap  = make(map[string]int)
	mapperGuard sync.RWMutex
)

func AddMapper(msgIns interface{}) {

	msgName := reflect.TypeOf(msgIns).String()

	msgID := cellnet.Name2ID(msgName)

	addMapper(msgName, msgID)
}

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
