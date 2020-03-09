package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/db"
	"ws/impl"
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
func main(){
	var wsPort= conf.GetConfig().WsPort
	http.HandleFunc("/", echo)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), nil);err!=nil{
		log.Fatal(err)
	}
}
func echo(w http.ResponseWriter, r *http.Request) {
	var (
		admin ResultAdmin
		total int
		err error
		wsConn *websocket.Conn
		conn *impl.Connection
		message []byte
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		goto Err
	}
	for {
		db.DB.First(&Admin{},1).Scan(&admin).Count(&total)
		log.Println("mysql id is:",total)
		if conn,err= impl.BuildConn(wsConn);err!=nil {
			goto Err
		}
		if message, err=conn.ReadMsg() ;err!=nil{
			log.Println("read:", err)
			goto  Err;
		}
		if err = conn.WriteMsg(message);err!=nil{
			log.Println("write:", err)
			goto Err
		}
		log.Printf("recv: %s", message)
	}
	Err:
		conn.Close()
}