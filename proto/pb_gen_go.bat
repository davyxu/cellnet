set outputpath=.
set protofile=%2%.proto
set packagename=%1%
set protoc_exe=protoc.exe
mkdir %outputpath%\%1%
"protoc.exe" %protofile% --plugin=protoc-gen-go=protoc-gen-go.exe --go_out %outputpath%\%packagename% --proto_path "."