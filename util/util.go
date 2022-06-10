//纯的工具类 永远用于被引入
package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
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

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal("zh_translations register failed")
	}
}

//判断路径 始终是以二进制所在路径为依据
func PathToEveryOne(path string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("路径错误")
	}
	baseDir := strings.Replace(dir, `\\`, "/", -1)
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

func translate(errs error) string {
	var errList []string
	for _, err := range errs.(validator.ValidationErrors) {
		errList = append(errList, err.Translate(trans))
	}
	return strings.Join(errList, ";")
}
func ValidateStruct(s interface{}) {
	errors := validate.Struct(s)
	if errors != nil {
		log.Fatal(translate(errors))
	}
}
func LogUtil(s, t string, debug bool) {
	if !debug {
		log.Println(s)
		return
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		_currentPath := PathToEveryOne("/")
		_file, _ := filepath.Abs(file)
		log.Printf("[%s] %s line=%d error is \n%s", t, strings.TrimLeft(_file, _currentPath), line, s)
	}

}

//防止野生goroutine panic 导致的整个程序退出
func Go(x func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Go run failed gorouting:", err)
		}
	}()
	x()
}
