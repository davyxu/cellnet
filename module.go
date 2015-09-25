package cellnet

import (
	"github.com/davyxu/cellnet/util"
	"log"
	"path"
	"runtime"
)

func getModuleName() string {

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return path.Dir(util.StripFileName(file, 3))
}

var moduleRegMap = make(map[string]func(Peer))

func RegisterModule(entry func(Peer)) {

	name := getModuleName()

	if GetModule(name) != nil {
		log.Println("duplicate module entry:", name)
		return
	}

	moduleRegMap[name] = entry
}

func GetModule(name string) func(Peer) {
	if v, ok := moduleRegMap[name]; ok {
		return v
	}

	return nil
}
