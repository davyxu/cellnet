package timer

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/example"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

func TestAfterTimer(t *testing.T) {

	signal := test.NewSignalTester(t)

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

func TestTickerTimer(t *testing.T) {

	signal := test.NewSignalTester(t)
	signal.SetTimeout(60 * time.Second)

	queue := cellnet.NewEventQueue()

	queue.StartLoop()

	var count int
	timer.Tick(queue, time.Millisecond*100, func(stopper timer.TickStopper) {
		log.Debugln("tick 100 ms", count)

		count++

		if count >= 10 {
			signal.Done(1)
			stopper.Stop()
		}

	})

	signal.WaitAndExpect("100ms * 10 times ticker not done", 1)
}

func TestDelay(t *testing.T) {

	signal := test.NewSignalTester(t)

	queue := cellnet.NewEventQueue()

	queue.StartLoop()

	log.Debugln("delay 1 sec begin")

	queue.DelayPost(time.Second, func() {

		log.Debugln("delay done")
		signal.Done(1)
	})

	signal.WaitAndExpect("delay not work", 1)
}
