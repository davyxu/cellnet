package nexus

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

func Register(disp *dispatcher.DataDispatcher) {

	if disp.Exists(cellnet.Type2ID(&coredef.ExpressACK{})) {
		panic("[nexus] Duplicate router register")
	}

	dispatcher.RegisterMessage(disp, coredef.ExpressACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.ExpressACK)

		if msg.GetCallID() != 0 {

			cellnet.InjectPost(cellnet.CellID(msg.GetTargetID()), &cellnet.Packet{
				MsgID: msg.GetMsgID(),
				Data:  msg.GetMsg(),
			}, msg.GetCallID())

		} else {

			cellnet.LocalPost(cellnet.CellID(msg.GetTargetID()), &cellnet.Packet{
				MsgID: msg.GetMsgID(),
				Data:  msg.GetMsg(),
			})
		}

	})

}

var (
	errExpressTargetNotFound error = errors.New("RPC reqest time out")
)

func init() {

	dispatcher.AddMapper(coredef.ExpressACK{})

	//注册节点系统的路由函数, 由addrlist来处理路由
	cellnet.SetExpressDriver(func(target cellnet.CellID, data interface{}, callid int64) error {
		// 取得目标所在的快递点信息
		rd := GetRegion(target.Region())
		if rd == nil {

			log.Println("[nexus] express target not found", target.String())
			return errExpressTargetNotFound
		}

		// 用户信息封包化
		pkt := data.(*cellnet.Packet)

		ack := &coredef.ExpressACK{
			Msg:      pkt.Data,
			MsgID:    proto.Uint32(pkt.MsgID),
			TargetID: proto.Int64(int64(target)),
		}

		if callid != 0 {
			ack.CallID = proto.Int64(callid)
		}

		// 先发到快递点, 再解包
		return cellnet.LocalPost(rd.Session, cellnet.BuildPacket(ack))
	})
}
