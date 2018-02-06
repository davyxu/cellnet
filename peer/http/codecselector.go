package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/httpform"
)

func selectCodec(method, contentType string) cellnet.Codec {
	if method == "GET" {
		return codec.MustGetCodec("httpform")
	}

	switch contentType {
	case "application/json":
		return codec.MustGetCodec("json")
	//case "application/xml", "text/xml":
	//	return XML
	default: //case "application/x-www-form-urlencoded", "multipart/form-data":
		return codec.MustGetCodec("httpform")
	}
}
