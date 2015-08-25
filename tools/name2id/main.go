package main

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"os"
)

func help() {
	fmt.Println("Usage: name2id NAME")
}

func main() {

	if len(os.Args) < 2 {
		help()
		return
	}

	id := cellnet.Name2ID(os.Args[1])

	fmt.Printf("id: %d (0x%x)", id, id)

}
