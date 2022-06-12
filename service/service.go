package service

// package service 为对外开发的接口 配合package router 对业务逻辑进行开发
func GetToken() string {
	return getToken()
}
func getToken() string {
	return "token"
}
