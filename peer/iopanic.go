package peer

type CoreCaptureIOPanic struct {
	captureIOPanic bool
}

func (self *CoreCaptureIOPanic) EnableCaptureIOPanic(v bool) {
	self.captureIOPanic = v
}

func (self *CoreCaptureIOPanic) CaptureIOPanic() bool {
	return self.captureIOPanic
}
