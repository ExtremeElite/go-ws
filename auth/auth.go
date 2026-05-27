package auth

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"ws/common"
	"ws/db"

	"github.com/golang-jwt/jwt"
)

var signKey []byte

type Claims struct {
	UserID      int
	ConnectType int
	jwt.StandardClaims
}

func init() {
	signKey = []byte(common.Conf.SignKey)
}

func CreateToken(id, connectType int) string {
	const maxAge = 86400
	c := &Claims{
		UserID:      id,
		ConnectType: connectType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(maxAge * time.Second).Unix(),
			Issuer:    common.Conf.Name,
			Id:        strconv.Itoa(id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := token.SignedString(signKey)
	if err != nil {
		log.Printf("[failed] create token: %v", err)
	}
	return s
}

func ParseToken(s string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(s, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func ValidateToken(token string, mold int) bool {
	switch mold {
	case 0:
		return true
	case 1:
		_, err := ParseToken(token)
		return err == nil
	case 2:
		return validateDB(token)
	default:
		return true
	}
}

func validateDB(token string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("validate token sql query:", r)
		}
	}()
	var total int
	if db.GormDB != nil {
		db.GormDB.Raw("SELECT COUNT(*) FROM admin WHERE id = ?", token).Scan(&total)
	}
	return total >= 1
}
