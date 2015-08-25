package nexus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

func Register(disp *dispatcher.PacketDispatcher) {

	if disp.Exists(cellnet.Type2ID(&coredef.RouterACK{})) {
		panic("[nexus] Duplicate router register")
	}

	dispatcher.RegisterMessage(disp, coredef.RouterACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.RouterACK)

		pkt := cellnet.Packet{
			MsgID: msg.GetMsgID(),
			Data:  msg.GetMsg(),
		}

		cellnet.SendLocal(cellnet.CellID(msg.GetTargetNodeID()), &pkt)

	})

}

func init() {

	//注册节点系统的路由函数, 由addrlist来处理路由
	cellnet.SetExpressDriver(func(target cellnet.CellID, data interface{}) bool {
		// 取得目标所在的快递点信息
		rd := GetRegion(target.Region())
		if rd == nil {
			log.Println("[nexus] express target not found", target.String())
			return false
		}

		// 用户信息封包化
		pkt := cellnet.BuildPacket(data.(proto.Message))

		// 先发到快递点, 再解包
		return cellnet.SendLocal(rd.Session, &coredef.RouterACK{
			Msg:          pkt.Data,
			MsgID:        proto.Uint32(pkt.MsgID),
			TargetNodeID: proto.Int64(int64(target)),
		})
	})
}
