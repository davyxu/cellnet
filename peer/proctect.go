package cellpeer

type Protect struct {
	CapturePanic bool
}

// 根据选项, 决定是否要捕获错误
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
