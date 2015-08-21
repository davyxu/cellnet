package cellnet

import (
	"log"
	"path"
	"runtime"
)

func getModuleName() string {

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return path.Dir(StripFileName(file, 3))
}

var moduleRegMap map[string]func() = make(map[string]func())

func RegisterModuleEntry(entry func()) {

	name := getModuleName()

	if GetModuleEntry(name) != nil {
		log.Println("duplicate module entry:", name)
		return
	}

	moduleRegMap[name] = entry
}

func GetModuleEntry(name string) func() {
	if v, ok := moduleRegMap[name]; ok {
		return v
	}

	return nil
}

func StartModule() {

	for name, entry := range moduleRegMap {
		log.Printf("start module: %v", name)

		entry()
	}
}
