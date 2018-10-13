package tests

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/timer"
)

func TestAfterTimer(t *testing.T) {

	signal := NewSignalTester(t)

	queue := cellnet.NewEventQueue()

	queue.StartLoop()

	timer.After(queue, 100*time.Millisecond, func() {
		log.Debugln("after 100 ms")

		signal.Done(1)
	}, nil)

	timer.After(queue, 200*time.Millisecond, func(context interface{}) {

		if context.(string) != "context" {
			t.FailNow()
		}

		log.Debugln("after 200 ms")

		signal.Done(2)
	}, "context")

	signal.WaitAndExpect("100ms after not done", 1)

	signal.WaitAndExpect("200ms after not done", 2)
}

func TestLoopTimer(t *testing.T) {

	signal := NewSignalTester(t)
	signal.SetTimeout(60 * time.Second)

	queue := cellnet.NewEventQueue()

	// 启动消息循环
	queue.StartLoop()

	var count int

	// 启动计时循环
	timer.NewLoop(queue, time.Millisecond*10, func(ctx *timer.Loop) {

		log.Debugln("tick 10 ms", count)

		count++

		if count >= 10 {
			signal.Done(1)
			ctx.Stop()
		}
	}, nil).Start()

	signal.WaitAndExpect("10ms * 10 times ticker not done", 1)
}
