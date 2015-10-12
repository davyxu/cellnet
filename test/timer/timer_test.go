package timer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/log"
	"github.com/davyxu/cellnet/test"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {

	signal := test.NewSignalTester(t)

	pipe := cellnet.NewEventPipe()

	evq := pipe.AddQueue()

	pipe.Start()

	const testTimes = 3

	var count int = testTimes

	cellnet.NewTimer(evq, time.Second, func(t *cellnet.Timer) {
		log.Debugln("timer 1 sec tick")

		signal.Done(1)

		count--

		if count == 0 {
			t.Stop()
			signal.Done(2)
		}
	})

	for i := 0; i < testTimes; i++ {
		signal.WaitAndExpect(1, "timer not tick")
	}

	signal.WaitAndExpect(2, "timer not stop")

}

func TestDelay(t *testing.T) {

	signal := test.NewSignalTester(t)

	pipe := cellnet.NewEventPipe()

	evq := pipe.AddQueue()

	pipe.Start()

	log.Debugln("delay 1 sec begin")

	evq.DelayPostData(time.Second, func() {

		log.Debugln("delay done")
		signal.Done(1)
	})

	signal.WaitAndExpect(1, "delay not work")
}
