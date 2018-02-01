package cellnet

type EventParam interface{}

type EventResult interface{}

// 事件函数的定义
type EventProc func(EventParam) EventResult
