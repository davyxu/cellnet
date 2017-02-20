cd pb
set PKGNAME=coredef
call gen_pb.bat coredef.proto
set PKGNAME=gamedef
call gen_pb.bat gamedef.proto

cd ..\sproto
call gen_sproto.bat core.sproto
cd ..