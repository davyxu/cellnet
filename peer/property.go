package peer

import (
	"reflect"
	"sync"
)

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

	switch rawValue := valuePtr.(type) {
	case *string:
		*rawValue = pv.(string)
	case *int:
		*rawValue = pv.(int)
	case *int32:
		*rawValue = pv.(int32)
	case *int64:
		*rawValue = pv.(int64)
	case *uint:
		*rawValue = pv.(uint)
	case *uint32:
		*rawValue = pv.(uint32)
	case *uint64:
		*rawValue = pv.(uint64)
	case *bool:
		*rawValue = pv.(bool)
	case *float32:
		*rawValue = pv.(float32)
	case *float64:
		*rawValue = pv.(float64)
	case *[]byte:
		*rawValue = pv.([]byte)
	default:
		v := reflect.Indirect(reflect.ValueOf(valuePtr))

		v.Set(reflect.ValueOf(pv))
	}

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
