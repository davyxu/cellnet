package socket

import (
	"fmt"
	"sync"
)

var (
	id2nameMap  = make(map[int]string)
	name2idMap  = make(map[string]int)
	mapperGuard sync.RWMutex
)

func MapNameID(name string, id int) {

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

	return fmt.Sprintf("(??)id=%d", id)
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
