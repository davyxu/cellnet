package cellcodec

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellnet"
	cellmeta "github.com/davyxu/cellnet/meta"
	"github.com/davyxu/x/container"
	xos "github.com/davyxu/x/os"
)

var (
	ErrNoCodec = errors.New("no codec")
)

// 编码消息, 在使用了带内存池的codec中，可以传入session或peer的ContextSet，保存内存池上下文，默认ctx传nil
func Encode(msg any, ps *xcontainer.Mapper) (data []byte, meta *cellmeta.Meta, err error) {

	// 获取消息元信息
	meta = cellmeta.MetaByMsg(msg)
	if meta == nil {
		return nil, nil, fmt.Errorf("msg not exists: %+v", msg)
	}

	if meta.Codec == nil {
		return nil, nil, ErrNoCodec
	}

	defer xos.Recover(func(raw any) {
		err = fmt.Errorf("encode panic: %+v", raw)
	})

	// 将消息编码为字节数组
	var raw any
	raw, err = meta.Codec.Encode(msg, ps)

	if err != nil {
		return
	}

	data = raw.([]byte)

	return
}

// 解码消息
func Decode(msgid int, data []byte) (msg any, meta *cellmeta.Meta, err error) {

	// 获取消息元信息
	meta = cellmeta.MetaById(msgid)

	// 消息没有注册
	if meta == nil {
		return nil, nil, fmt.Errorf("msg not exists: %d", msgid)
	}

	if meta.Codec == nil {
		return nil, nil, ErrNoCodec
	}

	defer xos.Recover(func(raw any) {
		err = fmt.Errorf("encode panic: %+v", raw)
	})

	// 创建消息
	msg = meta.NewType()

	// 从字节数组转换为消息
	err = meta.Codec.Decode(data, msg)

	return msg, meta, err
}

// Codec.Encode内分配的资源，在必要时可以回收，例如内存池对象
type CodecRecycler interface {
	Free(data any, ps *xcontainer.Mapper)
}

func Free(codec cellnet.Codec, data any, ps *xcontainer.Mapper) {

	if codec == nil {
		return
	}

	if recycler, ok := codec.(CodecRecycler); ok {
		recycler.Free(data, ps)
	}
}
