package http

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/httpform"
)

func selectCodec(method, contentType string) cellnet.Codec {
	if method == "GET" {
		return cellnet.MustGetCodec("httpform")
	}

	switch contentType {
	case "application/json":
		return cellnet.MustGetCodec("json")
	//case "application/xml", "text/xml":
	//	return XML
	default: //case "application/x-www-form-urlencoded", "multipart/form-data":
		return cellnet.MustGetCodec("httpform")
	}
}
