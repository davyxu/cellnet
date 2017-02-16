package cellnet

//// 消息到封包
//func BuildPacket(msgObj interface{}) ([]byte, error) {

//	msg := msgObj.(proto.Message)

//	rawdata, err := proto.Marshal(msg)

//	if err != nil {
//		return nil, err
//	}

//	return rawdata, nil
//}

//// 封包到消息
//func ParsePacket(data []byte, msgObj interface{}) error {
//	// msgType 为ptr类型, new时需要非ptr型

//	err := proto.Unmarshal(data, msgObj.(proto.Message))

//	if err != nil {
//		return err
//	}

//	return nil
//}
