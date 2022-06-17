# ws

#### 介绍

ws基础

#### 软件架构
- `Golang`实现的`Websocket`单体服务
  - 网络节点数据结构可配置:`Hashmap`|`Trie`
  - 数据库可配置：`Mysql`|`PostgreSQL`|`SQLite`|`Redis`
  - 鉴权认证可配置：`jwt验证`|`数据库验证`|`任意验证`
  - 消息体可配置:支持自定义消息体发送
  - 协议兼容：一统`Tcp`|`Udp`|`Websocket`|`Mqtt`|`Http`协议之间的消息传递
  - 并发安全：多线程+加锁的并发安全

- ws基础 win10 64位系统 内存不释放，32位正常释放 Linux 没测试;goland运行 并发未认证的情况会出现io wait 错误导致程序崩溃

###### port = 10090 web程序监听的端口

###### env = "dev" dev||prod dev打印日志 prod不打印
###### ws route /all 获取在线连接数

###### [mysqlDB]

###### serverHost = "127.0.0.1"

###### port = 3306

###### user = "root"

###### password = ""

###### db="nav"

###### maxConnect=5 数据库连接池维持的最大连接数量
###### e.g推送格式说明

```
{
    "event_type": 0, 
    "publish_account": [],
    "data": null,
}
event_type:数字事件1:转发事件;2:登录;3:退出登录;4:获取在线信息
publish_account:字符串数组代表推送的连接代表比如连接代表为[1,2,3]
data:被推送方需要收到的值可以是string、int 、object

```
## 配置config
- MultiplexPort 是否统一ws和http端口，如果统一那么，统一用ws端口对外暴露，对外暴露之后，ws端口发送请求必须通过认证,ws认证和http认证同样会被统一一个端口，一般通过数据库查询来验证token的合法性
- 数据库连接默认同时只支持一种数据库,目前只开发mysql数据库的支持
- 权限鉴权部分在可配置的情况下只支持一种权限验证jwt|数据库|任意验证
  
## 鉴权与认证
 + 通过验证jwt来对流模式进行鉴权，token的获取通过http post请求来获取。
    - http服务产生token，token会在长连接认证成功之后丢弃。
    - http请求是无状态的，可以对token的获取进行控制，可以获取多个token生产多个长链接。
    - 为了权限设计上的安全，token的使用为一次性的,鉴权成功之后，token会被丢弃,长连接不管什么原因导致的断开，如果需要再连接都需要重新通过http服务获取token。
    - 解耦长连接服务，http业务的开发远比长连接服务开发的工具链要多得多，http可以根据长连接的提供的接口开发对应的业务。诸如部分场景需要通过http的方式推送消息，让http业务层来对消息进行处理并且统一消息转发形式。为了更安全此服务的http的请求通过内网发送，不以任何形式对外暴露端口。