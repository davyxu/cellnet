set outdir=gamedef
set plugindir=..\..\..\..\..\bin
mkdir %outdir%
protoc.exe --plugin=protoc-gen-go=%plugindir%\protoc-gen-go.exe --go_out %outdir% --proto_path "." %*
protoc.exe --plugin=protoc-gen-msg=%plugindir%\protoc-gen-msg.exe --msg_out=msgid.go:%outdir% --proto_path "." %*