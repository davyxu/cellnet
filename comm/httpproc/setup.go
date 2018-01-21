package httpproc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"strings"
)

func ProcHttRequest(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.HttpEvent: // 接收数据事件

			contentType := ev.Req.Header.Get("Content-Type")

			if ev.Req.Method == "POST" || ev.Req.Method == "PUT" || ev.Req.Method == "PATCH" || contentType != "" {
				if strings.Contains(contentType, "form-urlencoded") {
					//context.Invoke(Form(obj, ifacePtr...))
				} else if strings.Contains(contentType, "multipart/form-data") {
					//context.Invoke(MultipartForm(obj, ifacePtr...))
				} else if strings.Contains(contentType, "json") {
					//context.Invoke(Json(obj, ifacePtr...))
				} else {
					//var errors Errors
					//if contentType == "" {
					//	errors.Add([]string{}, ContentTypeError, "Empty Content-Type")
					//} else {
					//	errors.Add([]string{}, ContentTypeError, "Unsupported Content-Type")
					//}
					//context.Map(errors)
				}
			} else {
				//context.Invoke(Form(obj, ifacePtr...))
			}

			//userFunc(&cellnet.RecvMsgEvent{ev.Ses, msg})
		default:
			userFunc(raw)
		}

		return nil
	}
}

func ProcForm(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.HttpEvent: // 接收数据事件

			//formStruct := reflect.New(reflect.TypeOf(formStruct))
			if err := ev.Req.ParseForm(); err != nil {
				return err
			}

			meta := cellnet.MessageMetaByName(ev.Req.URL.Path)
			if meta != nil {

				formStruct := reflect.ValueOf(meta.NewType())

				mapForm(formStruct, ev.Req.Form)

			}

		}

	}
}
