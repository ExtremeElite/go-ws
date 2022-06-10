# ws

#### 介绍

ws基础

#### 软件架构

ws基础 win10 64位系统 内存不释放，32位正常释放 Linux 没测试;goland运行 并发未认证的情况会出现io wait 错误导致程序崩溃

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
## config
- MultiplexPort 是否统一ws和http端口，如果统一那么，统一用ws端口对外暴露，对外暴露之后，ws端口发送请求必须通过认证,ws认证和http认证同样会被统一一个端口，一般通过数据库查询来验证token的合法性
  
## 鉴权与认证
 + 通过验证jwt来对流模式进行鉴权，token的获取通过http post请求来获取。
   
    - http服务产生token，token会在长连接认证成功之后丢弃。
    - http请求是无状态的，可以对token的获取进行控制，可以获取多个token生产多个长链接。
    - 为了权限设计上的安全，token的使用为一次性的,鉴权成功之后，token会被丢弃,长连接不管什么原因导致的断开，如果需要再连接都需要重新通过http服务获取token。