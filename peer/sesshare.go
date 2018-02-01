package peer

type CoreSessionIdentify struct {
	id int64
}

func (self *CoreSessionIdentify) ID() int64 {
	return self.id
}

func (self *CoreSessionIdentify) SetID(id int64) {
	self.id = id
}
