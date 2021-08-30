/*
date:2020\11\27 0027 22:03
email:gorouting@qq.com
author:gorouting
description:
*/
package broker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"ws/common"
	"ws/kernel"
)

func httpBroker(w http.ResponseWriter, r *http.Request) (err error) {
	var body []byte
	if body, err = validateData(w, r); err != nil {
		return
	}
	var pushData PushData
	if err = json.Unmarshal(body, &pushData); err != nil {
		return
	}
	workData(w, pushData)
	common.LogInfo("服务端收到:" + r.RemoteAddr + "发来的消息:" + string(body))
	return
}

//数据验证
func validateData(w http.ResponseWriter, r *http.Request) (body []byte, err error) {
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if bodyLen := len(body); bodyLen > common.Setting.MaxBody {
		res := `请求体大小为` + strconv.Itoa(bodyLen/1024) + `kb,大于` + strconv.Itoa(common.Setting.MaxBody/1024) + `kb`
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(res))
		common.LogInfo(res)
		return
	}
	return
}

//数据处理
func workData(w http.ResponseWriter, pushData PushData) {
	var response Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set(common.ContentType, common.AppJson)
	switch pushData.EventType {
	case Conversation:
		select {
		case HttpChan <- pushData:
		default:
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write(response.Json("ok", http.StatusOK, ""))
	case GetOnlineInfo:
		_, _ = w.Write(response.Json("ok", http.StatusOK, getOnLine()))
	default:
		_, _ = w.Write(response.Json("ok", http.StatusOK, pushData))
	}

}

//获取最新在线情况
func getOnLine() []string {
	result := make([]string, 0)
	nodes, _ := kernel.GetAllNode()
	for _, node := range nodes {
		result = append(result, node.Name)
	}
	return result
}

//转发http的数据到ws todo http消息转发一般为内网所以一般不需要进行身份认证
func HttpMessageForwarding() {
	for {
		select {
		case pushData := <-HttpChan:
			pushData.messageForwarding()
			common.LogInfo(`收到的http请求推送内容:` + pushData.ConversionJson())
		}
	}
}
