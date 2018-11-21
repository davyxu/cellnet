package peer

type CoreSQLParameter struct {
	PoolConnCount int
}

func (self *CoreSQLParameter) Init() {
	self.PoolConnCount = 1
}

func (self *CoreSQLParameter) SetPassword(v string) {
}

func (self *CoreSQLParameter) SetConnectionCount(v int) {
	self.PoolConnCount = v
}
