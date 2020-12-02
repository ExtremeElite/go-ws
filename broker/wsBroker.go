/**
 * @date:2020\11\27 0027 22:03
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	"ws/conf"
	"ws/core"
)

func wsBroker(conn *core.Connection) (err error){
	var message []byte
	if message, err=conn.ReadMsg();err!=nil{
		log.Println("读取失败:", err.Error())
		return
	}
	//读取消息并且发送消息
	if err=sendMessage(message,conn,
		ping,
		wsMessageForwarding,
	);err!=nil{
		return
	}
	log.Printf("服务器收到的: %s\n", message)
	return
}
func sendMessage(message []byte,conn *core.Connection,fns ...wsPipeLineFn) error  {
	for _,fn:=range fns {
		if err:=fn(message,conn);err!=nil{
			return err
		}
	}
	return nil
}
//客户端主动ping服务器
func ping(message []byte,conn *core.Connection) (err error) {
	if strings.ToLower(string(message))==`ping` {
		if err=conn.WriteMsg([]byte(`Pong`));err!=nil {
			log.Println("写入失败:", err.Error())
		}
		return
	}
	if strings.ToLower(string(message))==`pong` {
		if err=conn.WriteMsg([]byte(`Ping`));err!=nil {
			log.Println("写入失败:", err.Error())
		}
		return
	}
	return nil
}
//ws 消息转发 todo ws消息转发需要对连接权限进行认证
func wsMessageForwarding(message []byte,conn *core.Connection) (err error)  {
	var pushData PushData
	if err=json.Unmarshal(message,&pushData);err!=nil{
		log.Println("wsMessageForwarding err is:",err.Error())
		err=errors.New(`wsMessageForwarding data error`)
		return
	}
	switch pushData.EventType {
	case Conversation:
		pushData.messageForwarding()
	}
	return
}
//go程主动pong客户端
func pong(conn *core.Connection)  {
	var wsTimeOut=conf.CommonSet.WsTimeOut
	if wsTimeOut>0 {
		addTime:=time.Duration(wsTimeOut-1)*time.Second
		timer := time.NewTimer(addTime)
		for  {
			select {
			case <-timer.C:
				if conn.IsClose{
					timer.Stop()
					goto Over
				}
				conn.WriteMsg([]byte(`Pong`))
				log.Println(`Pong`)

			}
			timer.Reset(addTime)
		}
	}
	Over:
}