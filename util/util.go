package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (response Response) Json(msg string, code int, data interface{}) []byte {
	response.Msg = msg
	response.Code = code
	response.Data = data
	res, err := json.Marshal(response)
	if err != nil {
		return []byte(`["code":404,"msg":"数据错误","data":""]`)
	}
	return res
}

//判断路径 始终是以二进制所在路径为依据
func PathToEveryOne(path string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("路径错误")
	}
	baseDir := strings.Replace(dir, "\\", "/", -1)
	p, _ := filepath.Abs(baseDir + "/" + path)
	return p
}

//Unix时间戳转换为想要的格式
func Date(currentTime int64, currentDate ...string) string {
	layoutTime := "2006-01-02 15:04:05"
	if len(currentDate) > 0 {
		layoutTime = currentDate[0]
	}
	return time.Unix(currentTime, 0).Format(layoutTime)
}

// RandStringRunes 返回随机字符串
func RandString(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//返回两个数之间的随机数 左闭右闭
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	max = max + 1
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min) + min
}

//md5加密
func MD5(data string) string {
	h := md5.New()
	data = "QyAnrxYH7KGBJqMG4t0ymyVVJO5M2zgrP7bBjDL3LOM4PKJ8kOpzziuIrV0bcpXb" + data
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
