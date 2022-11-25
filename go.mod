module github.com/davyxu/cellnet

go 1.19

require (
	github.com/davyxu/protoplus v0.1.0
	github.com/davyxu/x v0.0.0
	github.com/davyxu/xlog v0.0.0
	github.com/golang/protobuf v1.5.2
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davyxu/ulexer v0.0.0-20200713054812-c9bb8db3521f // indirect
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/davyxu/protoplus => ../protoplus
	github.com/davyxu/x => ../x
	github.com/davyxu/xlog => ../xlog
)
