package jwt

import (
	sysLog "app/tools/logger"
	"app/tools/random"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	Key            = "admin123456!"
	DefaultExpTime = 7 * 86400 // jwt 默认过期时间（秒）
)

type UserJwt struct {
	UserId     uint   `json:"user_id"`
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
	jwt.StandardClaims
}

// CreateJwt 生成 jwt
func CreateJwt(id uint, expireTime int64) (string, *UserJwt) {
	if expireTime <= 0 {
		expireTime = time.Now().Unix() + DefaultExpTime
	}
	expireTime = time.Now().Unix() + expireTime

	userJwt := UserJwt{
		UserId:         id,
		Token:          random.Str(0),
		ExpireTime:     expireTime,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expireTime},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userJwt)
	tokenString, err := token.SignedString([]byte(Key))
	if err != nil {
		sysLog.Error("jwt 生成失败", err.Error())
		panic("jwt 生成失败" + err.Error())
	}
	return tokenString, &userJwt
}

func ParseJwt(token string) (*UserJwt, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserJwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Key), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserJwt); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
