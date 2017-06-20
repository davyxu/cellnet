package cellnet

type BranchHandler struct {
	hlist []EventHandler
}

func (self *BranchHandler) Call(ev *Event) {

	for _, h := range self.hlist {

		cloned := ev.Clone()

		HandlerLog(h, cloned)

		h.Call(cloned)

		if cloned.Result() != Result_OK {
			ev.SetResult(cloned.Result())
			break
		}
	}

}

func NewBranchHandler(hlist ...EventHandler) EventHandler {
	return &BranchHandler{
		hlist: hlist,
	}
}
