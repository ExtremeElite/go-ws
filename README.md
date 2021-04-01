# ws

#### 介绍

ws基础

#### 软件架构

ws基础 win10 64位系统 内存不释放，32位正常释放 Linux 没测试

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