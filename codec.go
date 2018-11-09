package cellnet

// 编码包
type Codec interface {
	// 将数据转换为字节数组
	Encode(msgObj interface{}, ctx ContextSet) (data interface{}, err error)

	// 将字节数组转换为数据
	Decode(data interface{}, msgObj interface{}) error

	// 编码器的名字
	Name() string

	// 兼容http类型
	MimeType() string
}
