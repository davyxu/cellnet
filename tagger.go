package cellnet

import "sync"

type tagData struct {
	name string
	data interface{}
}

type Tagger interface {
	GetTag(name string) (interface{}, bool)
	SetTag(name string, v interface{})
}

type CoreTagger struct {
	tag      []tagData
	tagGuard sync.RWMutex
}

func (self *CoreTagger) GetTag(name string) (interface{}, bool) {

	self.tagGuard.RLock()
	defer self.tagGuard.RUnlock()

	for _, t := range self.tag {
		if t.name == name {
			return t.data, true
		}
	}

	return nil, false
}

func (self *CoreTagger) SetTag(name string, v interface{}) {

	self.tagGuard.Lock()
	defer self.tagGuard.Unlock()

	for i, t := range self.tag {
		if t.name == name {
			self.tag[i] = tagData{name, v}
			return
		}
	}

	self.tag = append(self.tag, tagData{name, v})
}
