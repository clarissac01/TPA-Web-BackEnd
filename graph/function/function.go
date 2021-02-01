package function

import (
	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("mimi")

func ParseToken(tokenStr string) (float64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		return SecretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userid"].(float64)
		return userId, nil
	} else {
		return -1, err
	}
}


