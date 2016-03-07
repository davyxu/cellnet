package router

import (
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("router")

// 开启调试模式, 将显示完整的路由日志
var DebugMode bool
