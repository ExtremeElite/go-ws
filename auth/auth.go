package auth

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"ws/common"
	"ws/db"

	"github.com/golang-jwt/jwt/v5"
)

var signKey []byte

var (
	blacklist   = make(map[string]int64)
	blacklistMu sync.RWMutex
)

type Claims struct {
	UserID      int
	ConnectType int
	jwt.RegisteredClaims
}

func init() {
	signKey = []byte(common.Conf.SignKey)
}

func CreateToken(id, connectType int) string {
	const maxAge = 86400
	c := &Claims{
		UserID:      id,
		ConnectType: connectType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(maxAge * time.Second)),
			Issuer:    common.Conf.Name,
			ID:        strconv.Itoa(id),
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
		if err != nil {
			return false
		}
		blacklistMu.RLock()
		_, revoked := blacklist[token]
		blacklistMu.RUnlock()
		return !revoked
	case 2:
		return validateDB(token)
	default:
		return true
	}
}

func RevokeToken(token string) {
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	blacklist[token] = time.Now().Unix()
}

func IsRevoked(token string) bool {
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	_, exists := blacklist[token]
	return exists
}

func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			blacklistMu.Lock()
			now := time.Now().Unix()
			for t, exp := range blacklist {
				if exp < now {
					delete(blacklist, t)
				}
			}
			blacklistMu.Unlock()
		}
	}()
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
