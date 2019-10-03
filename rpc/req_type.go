package rpc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
	"time"
)

func CallType(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, userCallback interface{}) {
	callType(sesOrPeer, false, reqMsg, timeout, userCallback)
}

func CallSyncType(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, userCallback interface{}) {
	callType(sesOrPeer, true, reqMsg, timeout, userCallback)
}

// 异步RPC请求,按消息类型,一般用于客户端请求
// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType, error )
func callType(sesOrPeer interface{}, sync bool, reqMsg interface{}, timeout time.Duration, userCallback interface{}) {

	// 获取回调第一个参数
	funcType := reflect.TypeOf(userCallback)
	if funcType.Kind() != reflect.Func {
		panic("type rpc callback require 'func'")
	}

	if funcType.NumIn() != 2 {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType := funcType.In(0)
	if ackType.Kind() != reflect.Ptr {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType = ackType.Elem()

	callFunc := func(rawACK interface{}, err error) {
		vCall := reflect.ValueOf(userCallback)

		if rawACK == nil {
			rawACK = reflect.New(ackType).Interface()
		}

		var errV reflect.Value
		if err == nil {
			errV = nilError
		} else {
			errV = reflect.ValueOf(err)
		}

		vCall.Call([]reflect.Value{reflect.ValueOf(rawACK), errV})
	}

	ses, err := getPeerSession(sesOrPeer)

	if err != nil {
		callFunc(nil, err)
		return
	}

	createTypeRequest(sync, ackType, timeout, func() {
		ses.Send(reqMsg)
	}, callFunc)

}

var (
	nilError   = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
	callByType sync.Map // map[reflect.Type]func(interface{})
)

func createTypeRequest(sync bool, ackType reflect.Type, timeout time.Duration, onSend func(), onRecv func(rawACK interface{}, err error)) {

	if sync {
		feedBack := make(chan interface{})
		callByType.Store(ackType, feedBack)

		defer callByType.Delete(ackType)

		onSend()

		select {
		case ack := <-feedBack:
			onRecv(ack, nil)
		case <-time.After(timeout):
			onRecv(nil, ErrTimeout)
		}
	} else {

		callByType.Store(ackType, func(rawACK interface{}, err error) {
			onRecv(rawACK, err)
			callByType.Delete(ackType)
		})

		onSend()

		// 丢弃超时的类型,避免重复请求时,将第二次请求的消息删除
	}

}

type TypeRPCHooker struct {
}

func (TypeRPCHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	incomingMsgType := reflect.TypeOf(inputEvent.Message()).Elem()

	if rawFeedback, ok := callByType.Load(incomingMsgType); ok {

		switch feedBack := rawFeedback.(type) {
		case func(rawACK interface{}, err error):
			feedBack(inputEvent.Message(), nil)
		case chan interface{}:
			feedBack <- inputEvent.Message()
		}

		return inputEvent
	}

	return inputEvent
}

func (TypeRPCHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}
