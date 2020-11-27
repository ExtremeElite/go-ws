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
	var (
		body []byte
	)
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
	var pushData PushData
	if err=json.Unmarshal(body,&pushData);err!=nil{
		return
	}
	w.WriteHeader(http.StatusOK)
	if pushData.EventType== GetOnlineInfo {
		var response=Response{
			Code: 200,
			Msg:  "成功",
			Data: getOnLine(),
		}
		_, _ = w.Write(response.Json())
		return
	}
	select {
	case HttpChan <-pushData:
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	log.Println(r.RemoteAddr+"发来的消息:"+string(body))
	var response=Response{
		Code: 200,
		Msg:  "成功",
		Data: []string{},
	}
	_, _ = w.Write(response.Json())
	return
}
func getOnLine() []string {
	result :=make([]string,0)
	nodes,_:=core.GetAllNode()
	for _,node:= range nodes {
		result = append(result, node.Name)
	}
	return result
}