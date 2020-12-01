/**
 * @date:2020\11\27 0027 22:03
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/core"
)

func httpBroker(w http.ResponseWriter, r *http.Request)(err error)  {
	var body []byte
	if body,err=validateData(w,r);err!=nil{
		return
	}
	var pushData PushData
	if err=json.Unmarshal(body,&pushData);err!=nil{
		return
	}
	workData(w,pushData)
	log.Println(r.RemoteAddr+"发来的消息:"+string(body))
	return
}
//数据验证
func validateData(w http.ResponseWriter, r *http.Request) (body []byte, err error)  {
	if body,err=ioutil.ReadAll(r.Body);err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if bodyLen:=len(body);bodyLen>conf.CommonSet.MaxBody {
		res:=`请求体大小为`+strconv.Itoa(bodyLen/1024)+`kb,大于`+strconv.Itoa(conf.CommonSet.MaxBody/1024)+`kb`
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(res))
		log.Println(res)
		return
	}
	return
}
//数据处理
func workData(w http.ResponseWriter,pushData PushData){
	var response Response
	w.WriteHeader(http.StatusOK)
	switch pushData.EventType {
	case Conversation:
		select {
		case HttpChan <-pushData:
		default:
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write(response.Json("ok",http.StatusOK,""))
	case GetOnlineInfo:
		_, _ = w.Write(response.Json("ok",http.StatusOK,getOnLine()))
	default:
		_, _ = w.Write(response.Json("ok",http.StatusOK,pushData))
	}

}
//获取最新在线情况
func getOnLine() []string {
	result :=make([]string,0)
	nodes,_:=core.GetAllNode()
	for _,node:= range nodes {
		result = append(result, node.Name)
	}
	return result
}
//转发http的数据到ws
func GetDataFromHttp()  {
	for{
		select {
		case data:=<-HttpChan:
			core.Nodes.Range(func(name, node interface{}) bool {
				go func() {
					if len(data.PublishAccount)!=0 {
						for _,publishAccount:=range data.PublishAccount{
							if publishAccount==node.(*core.Node).Name {
								if err:=node.(*core.Node).Ws.WriteMsg([]byte(data.ConversionJson()));err!=nil{
									log.Println("data from http: ",err.Error())
								}
							}
						}
					}
				}()
				return true
			})
			log.Println(`收到的http请求推送内容:`+data.ConversionJson())
		}
	}
}