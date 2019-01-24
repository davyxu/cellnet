package codec

import (
	"fmt"
	"github.com/davyxu/cellnet"
)

var registedCodecs []cellnet.Codec

// 注册编码器
func RegisterCodec(c cellnet.Codec) {

	if GetCodec(c.Name()) != nil {
		panic("duplicate codec: " + c.Name())
	}

	registedCodecs = append(registedCodecs, c)
}

// 获取编码器
func GetCodec(name string) cellnet.Codec {

	for _, c := range registedCodecs {
		if c.Name() == name {
			return c
		}
	}

	return nil
}

// cellnet自带的编码对应包
func getPackageByCodecName(name string) string {
	switch name {
	case "binary":
		return "github.com/davyxu/cellnet/codec/binary"
	case "gogopb":
		return "github.com/davyxu/cellnet/codec/gogopb"
	case "httpjson":
		return "github.com/davyxu/cellnet/codec/httpjson"
	case "json":
		return "github.com/davyxu/cellnet/codec/json"
	case "protoplus":
		return "github.com/davyxu/cellnet/codec/protoplus"
	default:
		return "package/to/your/codec"
	}
}

// 指定编码器不存在时，报错
func MustGetCodec(name string) cellnet.Codec {
	codec := GetCodec(name)

	if codec == nil {
		panic(fmt.Sprintf("codec not found '%s'\ntry to add code below:\nimport (\n  _ \"%s\"\n)\n\n",
			name,
			getPackageByCodecName(name)))
	}

	return codec
}
