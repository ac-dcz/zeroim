package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenOption struct {
	AccessSecret  string //key
	AccessExpired int64  //过期时间
}

func BuildToken(opt *TokenOption, payLoad map[string]interface{}) (string, error) {
	iat := time.Now().Unix()
	exp := iat + opt.AccessExpired
	claims := make(jwt.MapClaims)
	claims["iat"] = iat
	claims["exp"] = exp
	for key, val := range payLoad {
		claims[key] = val
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(opt.AccessSecret))
}

func VerifyToken(opt *TokenOption, tokenString string) (map[string]interface{}, error) {
	token, err := jwt.NewParser(jwt.WithJSONNumber()).Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(opt.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if data, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("token claims error")
	} else {
		return data, nil
	}
}
