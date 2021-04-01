/**
 * @date:2021/4/1 16:13
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ws/core"
)

func AllNodeRouter() http.HandlerFunc {
	type Data struct {
		Info []core.Node `json:"info"`
		Total int `json:"total"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		nodes,total:=core.GetAllNode()
		for _,node:=range nodes {
			fmt.Printf("name:%s,remote_addr:%s\n",node.Name,node.RemoteAddr)
		}

		var data =Data{
			Info: nodes,
			Total: total,
		}
		result, err := json.Marshal(data)
		if err != nil {
			println(err.Error())
		}
		w.Write(result)
	}
}