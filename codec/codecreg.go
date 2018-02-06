package codec

import "github.com/davyxu/cellnet"

var registedCodecs []cellnet.Codec

func RegisterCodec(c cellnet.Codec) {

	if GetCodec(c.Name()) != nil {
		panic("duplicate codec: " + c.Name())
	}

	registedCodecs = append(registedCodecs, c)
}

func GetCodec(name string) cellnet.Codec {

	for _, c := range registedCodecs {
		if c.Name() == name {
			return c
		}
	}

	return nil
}

func MustGetCodec(name string) cellnet.Codec {
	codec := GetCodec(name)

	if codec == nil {
		panic("codec not register! " + name)
	}

	return codec
}
