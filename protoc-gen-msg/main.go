package main

import (
	"io/ioutil"
	"os"

	"bytes"
	"github.com/davyxu/golog"
	"github.com/davyxu/pbmeta"
	"github.com/gogo/protobuf/proto"
	pbprotos "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

var log = golog.New("main")

func main() {

	var errBuffer bytes.Buffer

	golog.SetOutput("main", &errBuffer)

	var Request plugin.CodeGeneratorRequest   // The input.
	var Response plugin.CodeGeneratorResponse // The output.

	defer func() {

		if errBuffer.Len() > 0 {
			str := errBuffer.String()
			Response.Error = &str
		}

		// 发回处理结果
		data, _ := proto.Marshal(&Response)

		os.Stdout.Write(data)

	}()

	// 读取protoc请求
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Errorln("reading input")
		os.Exit(1)
	}

	// 解析请求
	if err := proto.Unmarshal(data, &Request); err != nil {
		log.Errorln("parsing input proto")
		os.Exit(1)
	}

	if len(Request.FileToGenerate) == 0 {
		log.Errorln("no files to generate")
		os.Exit(1)
	}

	// 建立解析池
	pool := pbmeta.NewDescriptorPool(&pbprotos.FileDescriptorSet{
		File: Request.ProtoFile,
	})

	Response.File = make([]*plugin.CodeGeneratorResponse_File, 0)

	contenxt, ok := printFile(pool)

	if !ok {
		os.Exit(1)
	}

	Response.File = append(Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(Request.GetParameter()),
		Content: proto.String(contenxt),
	})

}
