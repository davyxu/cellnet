package cellnet

type Result int

const (
	Result_OK            Result = iota
	Result_SocketError          // 网络错误
	Result_SocketTimeout        // Socket超时
	Result_PackageCrack         // 封包破损
	Result_CodecError
)
