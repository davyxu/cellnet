package cellnet

import (
	"html/template"
	"net/http"
)

// Websocket接受器，具备会话访问
type WSAcceptor interface {
	GenericPeer

	SetHttps(certfile, keyfile string)

	// 设置升级器
	SetUpgrader(upgrader interface{})

	// 访问会话
	SessionAccessor
}

type HTTPAcceptor interface {
	GenericPeer

	// 设置http文件服务虚拟地址和文件系统根目录
	SetFileServe(dir string, root string)

	// 设置模板文件地址
	SetTemplateDir(dir string)

	// 设置http模板的分隔符，解决默认{{ }}冲突问题
	SetTemplateDelims(delimsLeft, delimsRight string)

	// 设置模板的扩展名，默认: .tpl .html
	SetTemplateExtensions(exts []string)

	// 设置模板函数入口
	SetTemplateFunc(f []template.FuncMap)
}

// HTTP连接器接口
type HTTPConnector interface {
	GenericPeer
	Request(method, path string, raw interface{}) (interface{}, error)
}

type HTTPSession interface {
	Request() *http.Request
}
