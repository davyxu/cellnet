# Protobuf编码安装

## Protobuf 编译器protoc下载

地址: https://github.com/google/protobuf/releases

说明：
* Mac OS 64位平台，选择 protoc-x.x.x-osx-x86_64.zip

* Windows 64位平台，选择 protoc-x.x.x-win32.zip

* Linux 64位平台，选择 protoc-x.x.x-linux-x86_64.zip

## protoc编译器

确保pb编译器protoc的可执行文件放在GOPATH/bin下


## 安装protoc编译器插件(protoc-gen-gogofaster)

该插件根据proto内容，生成xx.pb.go文件，包含pb的序列化相关代码

```
    go get -v github.com/gogo/protobuf/protoc-gen-gogofaster

    go install -v github.com/gogo/protobuf/protoc-gen-gogofaster
```

## 安装protoc编译器插件(protoc-gen-msg)

该插件根据proto内容，生成msgid.go文件，可以将消息绑定到cellnet的codec系统中，让cellnet可以识别protobuf的消息

```
    go get -v github.com/davyxu/cellnet/protoc-gen-msg

    go install -v github.com/davyxu/cellnet/protoc-gen-msg
```

## 测试

执行以下shell

```
${GOPATH}/github.com/davyxu/cellnet/codec/gogopb/test/export.sh
```

将使用protoc读取pb.proto并生成pb.pb.go和msgid.go两个文件

