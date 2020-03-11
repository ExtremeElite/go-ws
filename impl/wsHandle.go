package impl

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"ws/db"
)

type Admin struct {
	ID int `gorm:"primary_key"`
	UserName string `gorm:"column:user_name"`
	NickName string `gorm:"column:nick_name"`
	Status int8 `gorm:"column:status"`
	Password string `gorm:"column:password"`
	Ip string `gorm:"column:ip"`
	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

}
type ResultAdmin struct {
	ID uint `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nick_name"`
	Status int8 `json:"status"`
	HeaderImg string `json:"header_img"`
	UpdateTime int64 `json:"update_time"`
	Ip string `json:"ip"`
}
var upgrader=websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func WsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		admin ResultAdmin
		total int
		err error
		wsConn *websocket.Conn
		conn *Connection
		message []byte
	)
	if r.Header.Get("Connection")!="Upgrade" {
		w.Write([]byte("it is http"))
		return
	}
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		return
	}

	if conn,err= BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		goto Err
	}
	//登录判断
	if err=WsAuth(r);err!=nil {
		log.Println(err.Error())
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		goto Err
	}
	for {
		//超时设置
		conn.WsConn.SetReadDeadline(time.Now().Add(60*time.Second))
		db.DB.First(&Admin{},1).Scan(&admin).Count(&total)
		log.Println("mysql id is:",total)
		if message, err=conn.ReadMsg() ;err!=nil{
			log.Println("read:", err)
			goto Err
		}
		conn.Pong(message)
		log.Printf("recv: %s", message)
	}
Err:
	conn.Close()
}
