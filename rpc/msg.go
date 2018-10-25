package rpc

func (self *RemoteCallREQ) GetMsgID() uint16   { return uint16(self.MsgID) }
func (self *RemoteCallREQ) GetMsgData() []byte { return self.Data }
func (self *RemoteCallREQ) GetCallID() int64   { return self.CallID }
func (self *RemoteCallACK) GetMsgID() uint16   { return uint16(self.MsgID) }
func (self *RemoteCallACK) GetMsgData() []byte { return self.Data }
func (self *RemoteCallACK) GetCallID() int64   { return self.CallID }
