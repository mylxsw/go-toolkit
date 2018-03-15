package jwt_test

import (
	"testing"

	"git.yunsom.cn/golang/broadcast/utils/jwt"
)

var key = "123456"

func TestCreateToken(t *testing.T) {
	token := jwt.New(key)
	tokenString := token.CreateToken(jwt.Claims{
		"uid": 12344,
	}, -600)

	if tokenString == "" {
		t.Errorf("Create token failed")
	}
}

func TestParseToken(t *testing.T) {

	var uid int = 123456

	{
		// 检查token可用
		token := jwt.New(key)

		claims, err := token.ParseToken(token.CreateToken(jwt.Claims{
			"uid": uid,
		}, 500))
		if err != nil {
			t.Errorf("Parse token failed: %v", err)
		}

		if int(claims["uid"].(float64)) != uid {
			t.Errorf("Uid not match")
		}

	}
	// 不合法的token

	// 1. key错误
	{
		token := jwt.New(key)
		tokenStr := token.CreateToken(jwt.Claims{
			"uid": uid,
		}, 500)

		token2 := jwt.New("aabbcc")
		_, err := token2.ParseToken(tokenStr)
		if err == nil || err.Error() != "signature is invalid" {
			t.Errorf("Check invalid token failed")
		}
	}

	// 2. token 过期
	{
		token := jwt.New(key)
		tokenStr := token.CreateToken(jwt.Claims{
			"uid": uid,
		}, -5000)

		_, err := token.ParseToken(tokenStr)
		if err == nil || err.Error() != "Token is expired" {
			t.Errorf("Check invalid token failed")
		}
	}
}
