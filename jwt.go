package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ohdat/response"
	"github.com/spf13/viper"
	"time"
)

type JwtInfo struct {
	Uid        int    `json:"uid"`
	UnionID    string `json:"union_id"`
	OpenID     string `json:"open_id"`
	SessionKey string `json:"session_key"`
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	JwtInfo
}

//CreateToken 创建jwt Token
func CreateToken(info JwtInfo) (tokenString string, err error) {
	var secret = viper.GetString("jwt.secret")

	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    "OhdatLogin",
		},
		info,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secret))
	return
}

// ParseToken token解码
func ParseToken(tokenString string) (*JwtInfo, error) {
	var secret = viper.GetString("jwt.secret")

	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, response.ErrorTokenVerificationFail
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				//token 过期
				return nil, response.ErrorTokenExpire
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, response.ErrorTokenVerificationFail
			} else {
				return nil, response.ErrorTokenVerificationFail
			}

		}
	}
	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return &claims.JwtInfo, nil
	}
	return nil, fmt.Errorf("token无效")

}
