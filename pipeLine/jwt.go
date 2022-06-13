package pipeLine

import (
	"strconv"
	"time"
	"ws/common"

	"github.com/golang-jwt/jwt"
)

var SignKey = []byte(common.Setting.SignKey)

type CustomClaims struct {
	userId      int
	connectType int //1:ws 2:tcp 3:udp
	jwt.StandardClaims
}

func CreateToken(id, connectType int) string {
	maxAge := 60 * 60 * 24
	customClaims := &CustomClaims{
		userId:      id,
		connectType: connectType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(),
			Issuer:    common.Setting.Name,
			Id:        strconv.Itoa(id),
		},
	}
	//采用HMAC SHA256加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString(SignKey)
	if err != nil {
		common.LogInfoFailed(err.Error())
	}
	return tokenString
}
