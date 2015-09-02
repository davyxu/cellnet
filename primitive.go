package cellnet

import (
	"errors"
	"log"
	"sync"
)

var (
	// Cell实例管理
	cellMapGuard sync.RWMutex
	cellMap      map[CellID]*cell = make(map[CellID]*cell)

	// CellID生成器
	indexAccGuard sync.RWMutex
	indexAcc      int32
)

func genID() CellID {

	indexAccGuard.Lock()
	defer indexAccGuard.Unlock()

	// TODO 处理翻越case
	indexAcc++

	return NewCellID(RegionID, indexAcc)
}

func findCell(id CellID) *cell {
	cellMapGuard.RLock()
	defer cellMapGuard.RUnlock()

	if v, ok := cellMap[id]; ok {
		return v
	}

	return nil
}

// CellID是否为本进程内的ID
func IsLocal(id CellID) bool {
	return id.Region() == RegionID
}

// 为消息处理函数生成一个Cell, 返回CellID
func Spawn(callback func(interface{})) CellID {

	id := genID()

	if config.CellLog {
		log.Println("[cellnet] #spawn", id.String(), GetStackInfoString(2))
	}

	c := &cell{
		mailbox: make(chan interface{}, 8),
		id:      id,
	}

	cellMapGuard.Lock()
	cellMap[id] = c
	cellMapGuard.Unlock()

	go func() {

		for {

			if data, ok := c.fetch(); ok {
				callback(data)
			} else {
				break
			}

		}

	}()

	return id
}

var (
	errExpressDriverLost error = errors.New("Express driver lost")
	errTargetNotFound    error = errors.New("Target not found")
)

// 将制定内容发送到target的Cell中
func Send(target CellID, data interface{}) error {

	return RawSend(target, data, 0)
}

func RawSend(target CellID, data interface{}, callid int64) error {

	if target == 0 {
		return nil
	}

	if IsLocal(target) {
		return LocalPost(target, data)
	}

	return ExpressPost(target, data, callid)
}

// 将定制内容发送到远程
func ExpressPost(target CellID, data interface{}, callid int64) error {

	if expressDriver == nil {

		if config.CellLog {
			log.Println("[cellnet] express func nil, target not send", target.String())
		}

		return errExpressDriverLost
	}

	return expressDriver(target, data, callid)
}

// 将制定内容发送到本地的target的Cell中
func LocalPost(target CellID, data interface{}) error {

	if c := findCell(target); c != nil {

		if config.CellLog {
			log.Printf("[cellnet] #localpost %v %v %v", target.String(), ReflectContent(data), GetStackInfoString(3))
		}

		c.post(data)
		return nil
	}

	if config.CellLog {
		log.Println("[cellnet] target not found: ", target.String())
	}

	return errTargetNotFound
}

var expressDriver func(CellID, interface{}, int64) error

// 设置快递驱动, 负责将给定内容跨进程送达
func SetExpressDriver(driver func(CellID, interface{}, int64) error) {
	expressDriver = driver
}
