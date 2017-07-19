package timer

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

func TestAfterTimer(t *testing.T) {

	signal := util.NewSignalTester(t)

	queue := cellnet.NewEventQueue()

	queue.StartLoop()

	timer.After(queue, 500*time.Millisecond, func() {
		log.Debugln("after 100 ms")

		signal.Done(1)
	})

	timer.After(queue, 800*time.Millisecond, func() {
		log.Debugln("after 200 ms")

		signal.Done(2)
	})

	signal.WaitAndExpect("1 sec after not done", 1)

	signal.WaitAndExpect("2 sec after not done", 2)
}

func TestLoopTimer(t *testing.T) {

	signal := util.NewSignalTester(t)
	signal.SetTimeout(60 * time.Second)

	queue := cellnet.NewEventQueue()

	// 启动消息循环
	queue.StartLoop()

	var count int

	// 启动计时循环
	timer.NewLoop(queue, time.Millisecond*100, func(ctx *timer.Loop) {

		log.Debugln("tick 100 ms", count)

		count++

		if count >= 10 {
			signal.Done(1)
			ctx.Stop()
		}
	}, nil).Start()

	signal.WaitAndExpect("100ms * 10 times ticker not done", 1)
}
