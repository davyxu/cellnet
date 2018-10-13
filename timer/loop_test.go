package timer

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"testing"
	"time"
)

func TestLoopPanic(t *testing.T) {

	q := cellnet.NewEventQueue()
	q.EnableCapturePanic(true)

	q.StartLoop()

	var times = 3

	NewLoop(q, time.Millisecond*100, func(loop *Loop) {

		times--
		if times == 0 {
			loop.Stop()
			q.StopLoop()
		}

		fmt.Println("before")
		panic("panic")
		fmt.Println("after")

	}, nil).Start()

	q.Wait()
}
