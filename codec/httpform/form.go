package httpform

import (
	"github.com/davyxu/cellnet/codec"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type httpFormCodec struct {
}

const defaultMemory = 32 * 1024 * 1024

func (self *httpFormCodec) Name() string {
	return "httpform"
}

func (self *httpFormCodec) MimeType() string {
	return "application/x-www-form-urlencoded"
}

func structToUrlValues(obj interface{}) url.Values {
	objValue := reflect.Indirect(reflect.ValueOf(obj))

	objType := objValue.Type()

	var formValues = url.Values{}
	for i := 0; i < objValue.NumField(); i++ {

		fieldType := objType.Field(i)

		fieldValue := objValue.Field(i)

		formValues.Add(fieldType.Name, fieldValue.String())
	}

	return formValues
}

func (self *httpFormCodec) Encode(msgObj interface{}) (data interface{}, err error) {

	return strings.NewReader(structToUrlValues(msgObj).Encode()), err
}

func (self *httpFormCodec) Decode(data interface{}, msgObj interface{}) error {

	req := data.(*http.Request)

	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(defaultMemory)
	if err := mapForm(msgObj, req.Form); err != nil {
		return err
	}

	return nil
}

func init() {

	codec.RegisterCodec(new(httpFormCodec))
}
