package peer

import (
	"reflect"
	"sync"
)

type PropertySet interface {
	GetProperty(key, valuePtr interface{}) bool

	SetProperty(name interface{}, v interface{})

	RawGetProperty(name interface{}) (interface{}, bool)
}

type property struct {
	key   interface{}
	value interface{}
}

type CorePropertySet struct {
	properties      []property
	propertiesGuard sync.RWMutex
}

func (self *CorePropertySet) GetProperty(key, valuePtr interface{}) bool {

	pv, ok := self.RawGetProperty(key)
	if !ok {
		return false
	}

	v := reflect.Indirect(reflect.ValueOf(valuePtr))

	v.Set(reflect.ValueOf(pv))

	return true
}

func (self *CorePropertySet) RawGetProperty(key interface{}) (interface{}, bool) {

	self.propertiesGuard.RLock()
	defer self.propertiesGuard.RUnlock()

	for _, t := range self.properties {
		if t.key == key {
			return t.value, true
		}
	}

	return nil, false
}

func (self *CorePropertySet) SetProperty(key, v interface{}) {

	self.propertiesGuard.Lock()
	defer self.propertiesGuard.Unlock()

	for i, t := range self.properties {
		if t.key == key {
			self.properties[i] = property{key, v}
			return
		}
	}

	self.properties = append(self.properties, property{key, v})
}
