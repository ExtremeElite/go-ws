/**
 * @date:2021/4/1 16:13
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package router

import (
	"encoding/json"
	"net/http"
	"ws/common"
	"ws/kernel"
)

func AllNodeRouter() http.HandlerFunc {
	type Data struct {
		Info  []kernel.Node `json:"info"`
		Total int           `json:"total"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(common.ContentType, common.AppJson)
		nodes, total := kernel.GetAllNode()
		var data = Data{
			Info:  nodes,
			Total: total,
		}
		result, err := json.Marshal(data)
		if err != nil {
			println(err.Error())
			result = []byte(err.Error())
		}
		_, _ = w.Write(result)
	}
}
