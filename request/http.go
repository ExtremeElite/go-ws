/**
 * @date:2020/8/19 16:51
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package request

import (
	"net/http"
	"ws/middleware/auth"
)

func HttpRequest(w http.ResponseWriter,r *http.Request) (dataJson string,err error) {
	if r.Header.Get("Connection")!="Upgrade" {
		dataJson,err=auth.HttpAuth(r)
		if err!=nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(dataJson))
	}
	return
}
//ws地址的ws请求 包括数据验证
