/**
 * @date:2020\11\27 0027 22:03
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"log"
	"strings"
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

func ping(message []byte,conn *core.Connection) (err error) {
	if strings.ToLower(string(message))=="ping" {
		if err=conn.WriteMsg([]byte(`Pong`));err!=nil {
			log.Println("写入失败:", err.Error())
			return
		}
	}
	return nil
}
func pong(message []byte,conn *core.Connection) (err error) {
	if strings.ToLower(string(message))=="pong" {
		if err=conn.WriteMsg([]byte(`Ping`));err!=nil {
			log.Println("写入失败:", err.Error())
			return
		}
	}
	return nil
}

func wsMessageForwarding(message []byte,conn *core.Connection) (err error)  {
	var pushData PushData
	if err=json.Unmarshal(message,&pushData);err!=nil{
		return
	}
	switch pushData.EventType {
	case Conversation:
		pushData.messageForwarding()
	}
	return
}