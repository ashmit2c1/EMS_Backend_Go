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

func GetUserIDFromToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, checkfunction)
	if err != nil {
		return 0, errors.New("Could not parse the token")
	}
	if !parsedToken.Valid {
		return 0, errors.New("Token is not valid")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Could not extract claims from token")
	}
	userID, ok := claims["userID"].(float64) // Claims might be float64 in JWT
	if !ok {
		return 0, errors.New("userID not found in token")
	}

	return int64(userID), nil
}
