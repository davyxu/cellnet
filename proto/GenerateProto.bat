cd pb
set PKGNAME=gamedef
call gen_pb.bat gamedef.proto

cd ..\sproto
call gen_sproto.bat gamedef.sproto
cd ..

cd binary
call gen_binary
cd ..
