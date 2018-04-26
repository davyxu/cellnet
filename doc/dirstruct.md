# 目录功能

```
benchmark           性能测试

codec               编码支持，以及编码注册

    binary          二进制格式编码(github.com/davyxu/goobjfmt)

    httpform        http表单格式

    json            json编码格式

examples            例子

    chat            聊天

    echo            回音服务器

    fileserver      使用cellnet内建HTTP服务器支持文件服务

    websocket       WebSocket与网页js通信例子

msglog              消息日志处理

peer                各种协议的端实现，以及端注册入口及复用组件

    http            HTTP协议处理流程及端封装

    tcp             TCP协议处理流程及端封装

    udp             UDP协议处理流程及端封装

    gorillaws       WebSocket协议处理流程及端封装

proc                各种处理器实现，以及处理器注册入口

    http            HTTP消息处理及文件服务实现

    tcp             在TCP peer上构建的tcp处理器集合

    udp             在UDP peer上构建的udp处理器集合

    gorillaws       WeboScket的协议处理

relay               接力消息封装

rpc                 远程过程调用支持

tests               测试用例

timer               计时器接口

util                工具库

```