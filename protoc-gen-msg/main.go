package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/davyxu/golog"
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	plugin "github.com/davyxu/pbmeta/proto/compiler"
	"github.com/golang/protobuf/proto"
)

var log *golog.Logger = golog.New("main")

func main() {
	gen := New()

	// 读取protoc请求
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "reading input")
	}

	// 解析请求
	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err, "parsing input proto")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no files to generate")
	}

	// 建立解析池
	pool := pbmeta.NewDescriptorPool(&pbprotos.FileDescriptorSet{
		File: gen.Request.ProtoFile,
	})

	gen.Response.File = make([]*plugin.CodeGeneratorResponse_File, 0)

	for i := 0; i < pool.FileCount(); i++ {
		file := pool.File(i)

		gen.Reset()

		printFile(gen, file)

		gen.Response.File = append(gen.Response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(changeExt(file.FileName())),
			Content: proto.String(gen.String()),
		})

	}

	// 发回处理结果
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to write output proto")
	}

}

func changeExt(name string) string {
	ext := path.Ext(name)
	if ext == ".proto" {
		name = name[0 : len(name)-len(ext)]
	}
	return name + ".msg.go"
}
