package cellnet

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type context struct {
	name string
	data interface{}
}

// 消息元信息
type MessageMeta struct {
	Codec Codec        // 消息用到的编码
	Type  reflect.Type // 消息类型

	ID int // 消息ID (二进制协议中使用)

	ctxList []*context
}

func (self *MessageMeta) TypeName() string {

	if self == nil {
		return ""
	}

	if self.Type.Kind() == reflect.Ptr {
		return self.Type.Elem().Name()
	}

	return self.Type.Name()
}

func (self *MessageMeta) FullName() string {

	if self == nil {
		return ""
	}

	rtype := self.Type
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	var sb strings.Builder
	sb.WriteString(path.Base(rtype.PkgPath()))
	sb.WriteString(".")
	sb.WriteString(rtype.Name())

	return sb.String()
}

// 创建meta类型的实例
func (self *MessageMeta) NewType() interface{} {
	if self.Type == nil {
		return nil
	}

	return reflect.New(self.Type).Interface()
}

// 为meta对应的名字绑定上下文
func (self *MessageMeta) SetContext(name string, data interface{}) *MessageMeta {
	for _, ctx := range self.ctxList {

		if ctx.name == name {
			ctx.data = data
			return self
		}
	}

	self.ctxList = append(self.ctxList, &context{
		name: name,
		data: data,
	})

	return self
}

// 获取meta对应的名字绑定上下文
func (self *MessageMeta) GetContext(name string) (interface{}, bool) {

	for _, ctx := range self.ctxList {

		if ctx.name == name {
			return ctx.data, true
		}
	}

	return nil, false
}

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByFullName = map[string]*MessageMeta{}
	metaByID       = map[int]*MessageMeta{}
	metaByType     = map[reflect.Type]*MessageMeta{}
)

/*
http消息
Method URL -> Meta
Type -> Meta

非http消息
ID -> Meta
Type -> Meta

*/

// 注册消息元信息
func RegisterMessageMeta(meta *MessageMeta) *MessageMeta {

	// 非http类,才需要包装Type必须唯一

	if _, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by type: %d", meta.ID))
	} else {
		metaByType[meta.Type] = meta
	}

	if _, ok := metaByFullName[meta.FullName()]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by fullname: %s", meta.FullName()))
	} else {
		metaByFullName[meta.FullName()] = meta
	}

	if meta.ID == 0 {
		panic("message meta require 'ID' field: " + meta.TypeName())
	}

	if prev, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by id: %d type: %s, pre type: %s", meta.ID, meta.TypeName(), prev.TypeName()))
	} else {
		metaByID[meta.ID] = meta
	}

	return meta
}

// 根据名字查找消息元信息
func MessageMetaByFullName(name string) *MessageMeta {
	if v, ok := metaByFullName[name]; ok {
		return v
	}

	return nil
}

func MessageMetaVisit(nameRule string, callback func(meta *MessageMeta) bool) error {
	exp, err := regexp.Compile(nameRule)
	if err != nil {
		return err
	}

	for name, meta := range metaByFullName {
		if exp.MatchString(name) {

			if !callback(meta) {
				return nil
			}

		}
	}

	return nil
}

// 根据类型查找消息元信息
func MessageMetaByType(t reflect.Type) *MessageMeta {

	if t == nil {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByType[t]; ok {
		return v
	}

	return nil
}

// 根据消息对象获得消息元信息
func MessageMetaByMsg(msg interface{}) *MessageMeta {

	if msg == nil {
		return nil
	}

	return MessageMetaByType(reflect.TypeOf(msg))
}

// 根据id查找消息元信息
func MessageMetaByID(id int) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

func MessageToName(msg interface{}) string {

	if msg == nil {
		return ""
	}

	meta := MessageMetaByMsg(msg)
	if meta == nil {
		return ""
	}

	return meta.TypeName()
}

func MessageToID(msg interface{}) int {

	if msg == nil {
		return 0
	}

	meta := MessageMetaByMsg(msg)
	if meta == nil {
		return 0
	}

	return int(meta.ID)
}

func MessageSize(msg interface{}) int {

	if msg == nil {
		return 0
	}

	// 获取消息元信息
	meta := MessageMetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return 0
	}

	// 将消息编码为字节数组
	raw, err := meta.Codec.Encode(msg, nil)

	if err != nil {
		return 0
	}

	return len(raw.([]byte))
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
