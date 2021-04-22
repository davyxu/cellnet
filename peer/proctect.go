package cellpeer

type Protect struct {
	CapturePanic bool
}

func (self *Protect) ProctectCall(job func(), cleanup func(raw interface{})) {

	if self.CapturePanic {

		defer func() {
			if raw := recover(); raw != nil {
				cleanup(raw)
			}
		}()

		job()
	} else {
		job()
	}
}
