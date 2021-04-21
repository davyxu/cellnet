package cellmeta

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByFullName = map[string]*Meta{}
	metaByID       = map[int]*Meta{}
	metaByType     = map[reflect.Type]*Meta{}
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
func Register(meta *Meta) *Meta {

	// 注册时, 统一为非指针类型
	if meta.Type.Kind() == reflect.Ptr {
		meta.Type = meta.Type.Elem()
	}

	meta.init()

	if pre, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by id: %d name: %s, pre id: %d name: %s", meta.ID, meta.Type.Name(), pre.ID, pre.TypeName()))
	} else {
		metaByType[meta.Type] = meta
	}

	if _, ok := metaByFullName[meta.FullName()]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by fullname: %s", meta.FullName()))
	} else {
		metaByFullName[meta.FullName()] = meta
	}

	if meta.ID != 0 {
		if prev, ok := metaByID[meta.ID]; ok {
			panic(fmt.Sprintf("Duplicate message meta register by id: %d type: %s, pre type: %s", meta.ID, meta.TypeName(), prev.TypeName()))
		} else {
			metaByID[meta.ID] = meta
		}
	}

	return meta
}

// 根据名字查找消息元信息
func MetaByFullName(name string) *Meta {
	if v, ok := metaByFullName[name]; ok {
		return v
	}

	return nil
}

func MetaVisit(nameRule string, callback func(meta *Meta) bool) error {
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
func MetaByType(t reflect.Type) *Meta {

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
func MetaByMsg(msg interface{}) *Meta {

	if msg == nil {
		return nil
	}

	return MetaByType(reflect.TypeOf(msg))
}

// 根据id查找消息元信息
func MetaByID(id int) *Meta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

// 消息名(例如:MsgREQ)
func MessageToName(msg interface{}) string {

	if msg == nil {
		return ""
	}

	meta := MetaByMsg(msg)
	if meta == nil {
		return ""
	}

	return meta.TypeName()
}

func MessageIDToName(msgid int) string {
	meta := MetaByID(msgid)
	if meta != nil {
		return meta.TypeName()
	}

	return ""
}

// 消息名(例如:proto.MsgREQ)
func MessageToFullName(msg interface{}) string {

	if msg == nil {
		return ""
	}

	meta := MetaByMsg(msg)
	if meta == nil {
		return ""
	}

	return meta.FullName()
}

func MessageToID(msg interface{}) int {

	if msg == nil {
		return 0
	}

	meta := MetaByMsg(msg)
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
	meta := MetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return 0
	}

	if meta.Codec == nil {
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
