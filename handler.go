package cellnet

type Handler interface {
	Call(int, interface{}) error
}

//func (self *Handler) Call(ev int, data interface{}, h **Handler) error {

//	*h = self.next

//	return self.entry(ev, data)
//}

//func LinkHandler(list ...func(int, interface{}) error) (head *Handler) {

//	var prev **Handler
//	for _, v := range list {

//		this := &Handler{
//			entry: v,
//		}

//		if prev != nil {
//			(*prev).next = this

//		} else {
//			head = this
//		}

//		prev = &this

//	}

//	return head

//}
