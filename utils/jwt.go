package utils

import (
	"errors"
	"time"

	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString(secretKey)
}

func checkfunction(token *jwt.Token) (interface{}, error) {
	_, err := token.Method.(*jwt.SigningMethodHMAC)
	if err == false {
		return nil, errors.New("Unexpected Signing Method")
	}
	return secretKey, nil
}
func VerifyToken(token string) error {
	parsedToken, err := jwt.Parse(token, checkfunction)
	if err != nil {
		return errors.New("Could not parse the token")
	}
	check := parsedToken.Valid
	if check == false {
		return errors.New("Token is not valid")
	}
	return nil
}
