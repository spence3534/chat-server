package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotVaildYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

var mySigningKey = []byte("chatServerKey")

func GetToken(name string) string {

	c := TokenClaims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			// 生效时间
			NotBefore: time.Now().Unix() - 60,
			// 过期时间 这里是两天
			ExpiresAt: time.Now().Unix() + 60*60*24*2,
			// 签发人
			Issuer: "administrator",
		},
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	s, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("生成token失败", err)
	}
	return s
}

// token解密

func ParseToken(tokenStr string) (*TokenClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	fmt.Println(err)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// 不是token
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// token过期
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// 错误的token
				return nil, TokenNotVaildYet
			} else {
				return nil, TokenInvalid
			}

		}
	}

	if t != nil {
		if claims, ok := t.Claims.(*TokenClaims); ok {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
