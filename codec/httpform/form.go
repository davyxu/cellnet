package httpform

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

func anyToString(any interface{}) string {

	switch v := any.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		panic("Unknown type to convert to string")
	}
}

func structToUrlValues(obj interface{}) url.Values {
	objValue := reflect.Indirect(reflect.ValueOf(obj))

	objType := objValue.Type()

	var formValues = url.Values{}
	for i := 0; i < objValue.NumField(); i++ {

		fieldType := objType.Field(i)

		fieldValue := objValue.Field(i)

		formValues.Add(fieldType.Name, anyToString(fieldValue.Interface()))
	}

	return formValues
}

func (self *httpFormCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

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
