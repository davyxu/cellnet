/*
nexus包为cellnet提供了跨进程,机器的访问支持
其基本设计的思想是去中心化, 不会因为中心宕机而导致功能受限

每个独立操作系统进程就是一个region, region间根据启动参数指定连接点

任何一个region只要连接上其他region,就会通过addresslist共享方式获得网络内所有
其他region的连接信息同时连接过去

nexus只是跨进程互联提供的一种驱动,未来可以根据需要提供更多互联方式

*/
package nexus
