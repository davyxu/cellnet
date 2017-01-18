package cellnet

import (
	"path"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   uint32
}

var (
	name2msgmeta = map[string]*MessageMeta{}
	id2msgmeta   = map[uint32]*MessageMeta{}
)

// 注册消息元信息(代码生成专用)
func RegisterMessageMeta(name string, msg proto.Message, id uint32) {

	rtype := reflect.TypeOf(msg)

	meta := &MessageMeta{
		Type: rtype,
		Name: name,
		ID:   id,
	}

	name2msgmeta[name] = meta
	id2msgmeta[meta.ID] = meta
}

// 根据名字查找消息元信息
func MessageMetaByName(name string) *MessageMeta {
	if v, ok := name2msgmeta[name]; ok {
		return v
	}

	return nil
}

func MessageFullName(rtype reflect.Type) string {

	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	return path.Base(rtype.PkgPath()) + "." + rtype.Name()

}

// 根据id查找消息元信息
func MessageMetaByID(id uint32) *MessageMeta {
	if v, ok := id2msgmeta[id]; ok {
		return v
	}

	return nil
}

// 根据id查找消息名, 没找到返回空
func MessageNameByID(id uint32) string {

	if meta := MessageMetaByID(id); meta != nil {
		return meta.Name
	}

	return ""
}

// 遍历消息元信息
func VisitMessageMeta(callback func(*MessageMeta)) {

	for _, meta := range name2msgmeta {
		callback(meta)
	}

}
