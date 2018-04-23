package codec

import "github.com/davyxu/cellnet"

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

// 指定编码器不存在时，报错
func MustGetCodec(name string) cellnet.Codec {
	codec := GetCodec(name)

	if codec == nil {
		panic("codec not register! " + name)
	}

	return codec
}
