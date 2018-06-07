package cellnet

import (
	"fmt"
	"reflect"
)

/*
客户端请求时  Method+RequestType -> Meta

服务器响应请求 URL+Method -> Meta

服务器回应 Method+ResponseType -> Meta

客户端接收 URL+Method -> Meta

*/

// 消息元信息
type HttpMeta struct {
	Path   string // 相对路径
	Method string // 方式

	RequestCodec Codec        // 请求消息编码
	RequestType  reflect.Type // 请求消息类型

	ResponseCodec Codec        // 响应消息编码
	ResponseType  reflect.Type // 响应消息类型
}

func (self *HttpMeta) RequestTypeName() string {

	if self == nil {
		return ""
	}

	if self.RequestType.Kind() == reflect.Ptr {
		return self.RequestType.Elem().Name()
	}

	return self.RequestType.Name()
}

func (self *HttpMeta) ResponseTypeName() string {

	if self == nil {
		return ""
	}

	if self.ResponseType.Kind() == reflect.Ptr {
		return self.ResponseType.Elem().Name()
	}

	return self.ResponseType.Name()
}

type methodURLPair struct {
	Method string
	URL    string
}

type methodTypePair struct {
	Method string
	reflect.Type
}

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByHttpPair     = map[methodURLPair]*HttpMeta{}
	metaByRequestType  = map[methodTypePair]*HttpMeta{}
	metaByResponseType = map[methodTypePair]*HttpMeta{}
)

func RegisterHttpMeta(meta *HttpMeta) {

	urlPair := methodURLPair{meta.Method, meta.Path}

	reqPair := methodTypePair{meta.Method, meta.RequestType}

	respPair := methodTypePair{meta.Method, meta.ResponseType}

	if _, ok := metaByHttpPair[urlPair]; ok {
		panic("Duplicate message meta register by URL: " + meta.Path)
	} else {
		metaByHttpPair[urlPair] = meta
	}

	if _, ok := metaByRequestType[reqPair]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by request type: %s", meta.Path))
	} else {
		metaByRequestType[reqPair] = meta
	}

	if _, ok := metaByResponseType[respPair]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by response type: %s", meta.Path))
	} else {
		metaByResponseType[respPair] = meta
	}

}

func HttpMetaByMethodURL(method, url string) *HttpMeta {
	if v, ok := metaByHttpPair[methodURLPair{method, url}]; ok {
		return v
	}

	return nil
}

func HttpMetaByRequestType(method string, t reflect.Type) *HttpMeta {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByRequestType[methodTypePair{method, t}]; ok {
		return v
	}

	return nil
}

func HttpMetaByResponseType(method string, t reflect.Type) *HttpMeta {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByResponseType[methodTypePair{method, t}]; ok {
		return v
	}

	return nil
}
