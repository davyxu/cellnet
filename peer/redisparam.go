package peer

type CoreRedisParameter struct {
	Password      string
	DBIndex       int
	PoolConnCount int
}

func (self *CoreRedisParameter) Init() {
	self.PoolConnCount = 1
}

func (self *CoreRedisParameter) SetPassword(v string) {
	self.Password = v
}

func (self *CoreRedisParameter) SetDBIndex(v int) {
	self.DBIndex = v
}

func (self *CoreRedisParameter) SetConnectionCount(v int) {
	self.PoolConnCount = v
}
