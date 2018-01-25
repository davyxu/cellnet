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

func MakeEventProcessor(name string, inbound, outbound EventProc) (EventProc, EventProc) {

	f := evtprocByName[name]
	if f == nil {
		panic("Event processor not found: " + name)
	}

	return f(inbound, outbound)
}

type DuplexEventInitor interface {
	SetEventFunc(processor string, inboundEvent, outboundEvent EventProc)
	SetRaw(inboundEvent, outboundEvent EventProc)
}

type DuplexEventInvoker interface {
	CallInboundProc(ev interface{}) interface{}

	CallOutboundProc(ev interface{}) interface{}
}

type CoreDuplexEventProc struct {
	InboundProc  EventProc
	OutboundProc EventProc
}

func (self *CoreDuplexEventProc) SetEventFunc(processor string, inboundEvent, outboundEvent EventProc) {
	self.InboundProc, self.OutboundProc = MakeEventProcessor(processor, inboundEvent, outboundEvent)
}

func (self *CoreDuplexEventProc) SetRaw(inboundEvent, outboundEvent EventProc) {
	self.InboundProc, self.OutboundProc = inboundEvent, outboundEvent
}

// socket包内部派发事件
func (self *CoreDuplexEventProc) CallInboundProc(ev interface{}) interface{} {

	if self.InboundProc == nil {
		return nil
	}

	//log.Debugf("<Inbound> %T|%+v", ev, ev)

	return self.InboundProc(ev)
}

// socket包内部派发事件
func (self *CoreDuplexEventProc) CallOutboundProc(ev interface{}) interface{} {

	if self.OutboundProc == nil {
		return nil
	}

	//log.Debugf("<Outbound> %T|%+v", ev, ev)

	return self.OutboundProc(ev)
}
