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

func (self *sprotoCodec) Encode(msgObj interface{}) ([]byte, error) {

	return sproto.Encode(msgObj)
}

func (self *sprotoCodec) Decode(data []byte, msgObj interface{}) error {

	// sproto要求必须有头, 但空包也是可以的
	if len(data) == 0 {
		return nil
	}

	_, err := sproto.Decode(data, msgObj)

	return err
}

func AutoRegisterMessageMeta(msgTypes []reflect.Type) {

	for _, tp := range msgTypes {

		msgName := fmt.Sprintf("%s.%s", path.Base(tp.PkgPath()), tp.Name())

		cellnet.RegisterMessageMeta(msgName, tp, util.StringHash(msgName))
	}

}

func init() {

	cellnet.RegisterCodec("sproto", new(sprotoCodec))
}
