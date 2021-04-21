package cellcodec

import (
	"fmt"
	"github.com/davyxu/cellnet"
)

var (
	codecByName = map[string]cellnet5.Codec{}
)

// 注册编码器
func Register(c cellnet5.Codec) {

	if codecByName[c.Name()] != nil {
		panic("duplicate codec: " + c.Name())
	}

	codecByName[c.Name()] = c
}

// 获取编码器
func GetByName(name string) cellnet5.Codec {

	for _, c := range codecByName {
		if c.Name() == name {
			return c
		}
	}

	return nil
}

// cellnet自带的编码对应包
func getPackageByCodecName(name string) string {
	switch name {
	case "gogopb":
		return "github.com/davyxu/cellnet/codec/gogopb"
	case "json":
		return "github.com/davyxu/cellnet/codec/json"
	case "protoplus":
		return "github.com/davyxu/cellnet/codec/protoplus"
	default:
		return "package/to/your/codec"
	}
}

// 指定编码器不存在时，报错
func MustGetByName(name string) cellnet5.Codec {
	codec := GetByName(name)

	if codec == nil {
		panic(fmt.Sprintf("codec not found '%s'\ntry to add code below:\nimport (\n  _ \"%s\"\n)\n\n",
			name,
			getPackageByCodecName(name)))
	}

	return codec
}
