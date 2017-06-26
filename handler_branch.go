package cellnet

import "sync"

type BranchHandler struct {
	hlist      []EventHandler
	hlistGuard sync.RWMutex
}

func (self *BranchHandler) Call(ev *Event) {

	for _, h := range self.List() {

		cloned := ev.Clone()

		HandlerLog(h, cloned)

		h.Call(cloned)

		if cloned.Result() != Result_OK {
			ev.SetResult(cloned.Result())
			break
		}
	}

}

func (self *BranchHandler) List() []EventHandler {
	self.hlistGuard.RLock()
	defer self.hlistGuard.RUnlock()

	return self.hlist
}

func (self *BranchHandler) Add(h EventHandler) {
	self.hlistGuard.Lock()
	self.hlist = append(self.hlist, h)
	self.hlistGuard.Unlock()
}

func NewBranchHandler(hlist ...EventHandler) EventHandler {
	return &BranchHandler{
		hlist: hlist,
	}
}
