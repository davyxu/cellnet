package sprotocodec

import (
	"fmt"
	"path"
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/gosproto"
)

type sprotoCodec struct {
}

func (self *sprotoCodec) Name() string {
	return "sproto"
}

func (self *sprotoCodec) Encode(msgObj interface{}) ([]byte, error) {

	data, err := sproto.Encode(msgObj)
	if err != nil {
		return nil, err
	}

	return sproto.Pack(data), nil
}

func (self *sprotoCodec) Decode(data []byte, msgObj interface{}) error {

	// sproto要求必须有头, 但空包也是可以的
	if len(data) == 0 {
		return nil
	}

	raw, err := sproto.Unpack(data)
	if err != nil {
		return err
	}

	_, err2 := sproto.Decode(raw, msgObj)

	return err2
}

func AutoRegisterMessageMeta(msgTypes []reflect.Type) {

	for _, tp := range msgTypes {

		msgName := fmt.Sprintf("%s.%s", path.Base(tp.PkgPath()), tp.Name())

		cellnet.RegisterMessageMeta("sproto", msgName, tp, util.StringHash(msgName))
	}

}

func init() {

	cellnet.RegisterCodec("sproto", new(sprotoCodec))
}
