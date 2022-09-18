module github.com/davyxu/cellnet

go 1.15

require (
	github.com/davyxu/protoplus v0.1.0
	github.com/davyxu/x v0.0.0
	github.com/davyxu/xlog v0.0.0
	github.com/golang/protobuf v1.5.2
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/davyxu/protoplus => ../protoplus
	github.com/davyxu/x => ../x
	github.com/davyxu/xlog => ../xlog
)
