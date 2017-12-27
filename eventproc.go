package cellnet

type EventParam interface{}

type EventResult interface{}

// 事件函数的定义
type EventProc func(EventParam) EventResult

type DuplexEventProcessor func(EventProc, EventProc) (EventProc, EventProc)

var evtprocByName = map[string]DuplexEventProcessor{}

func RegisterEventProcessor(name string, f DuplexEventProcessor) {

	if _, ok := creatorByTypeName[name]; ok {
		panic("Duplicate peer type")
	}

	evtprocByName[name] = f
}

func GetEventProcessor(name string, inbound, outbound EventProc) (EventProc, EventProc) {

	f := evtprocByName[name]
	if f == nil {
		panic("Event processor not found: " + name)
	}

	return f(inbound, outbound)
}
