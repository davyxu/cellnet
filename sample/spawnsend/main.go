package main

import (
	"github.com/davyxu/cellnet"
	"log"
	"time"
)

var done = make(chan bool)

func spawnsend() {

	// no block spawn cell, msg function here
	cid := cellnet.Spawn(func(mailbox chan interface{}) {
		for {

			switch v := (<-mailbox).(type) {
			case string:
				log.Println(v)
			}
		}

	})

	cellnet.Send(cid, "hello world ")

	done <- true
}

func main() {

	go spawnsend()

	select {
	case <-done:

	case <-time.After(3 * time.Second):
		log.Println("time out")
	}

}
