package cellnet

import (
	"reflect"

	"github.com/davyxu/cellnet/util"
	"github.com/golang/protobuf/proto"
)

type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   int
}

var (
	name2msgmeta    = make(map[string]*MessageMeta)
	id2msgmeta      = make(map[int]*MessageMeta)
	msgtype2msgmeta = make(map[reflect.Type]*MessageMeta)
)

// 注册消息元信息(代码生成专用)
func RegisterMessageMeta(name string, msg proto.Message) {

	rtype := reflect.TypeOf(msg)

	meta := &MessageMeta{
		Type: rtype,
		Name: name,
		ID:   int(util.StringHashNoCase(name)),
	}

	name2msgmeta[name] = meta
	msgtype2msgmeta[rtype] = meta
	id2msgmeta[meta.ID] = meta
}

// 根据名字查找消息元信息
func MessageMetaByName(name string) *MessageMeta {
	if v, ok := name2msgmeta[name]; ok {
		return v
	}

	return nil
}

// 根据类型名字查找消息元信息
func MessageMetaByType(rtype reflect.Type) *MessageMeta {
	if v, ok := msgtype2msgmeta[rtype]; ok {
		return v
	}

	return nil
}

// 根据id查找消息元信息
func MessageMetaByID(id int) *MessageMeta {
	if v, ok := id2msgmeta[id]; ok {
		return v
	}

	return nil
}

// 遍历消息元信息
func VisitMessageMeta(callback func(*MessageMeta)) {

	for _, meta := range name2msgmeta {
		callback(meta)
	}

}
