module github.com/davyxu/cellnet

go 1.15

require (
	github.com/davyxu/protoplus v0.1.0
	github.com/davyxu/x v0.0.0
	github.com/golang/protobuf v1.5.0
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/davyxu/protoplus => ../protoplus

replace github.com/davyxu/x => ../x
