package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken(id string, username string, exp int32) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
		"id": id,
        "username": username,
        "exp": time.Now().Add(time.Second * time.Duration(exp)).Unix(),
        })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
    if err != nil {
    return "", err
    }

 return tokenString, nil
}

func CreateAccessToken(id string, username string, exp int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
		"id": id,
		"username": username,
		"exp": time.Now().Add(time.Second * time.Duration(exp)).Unix(),
		})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {

	return "", err
	}

 return tokenString, nil
}

func VerifyRefreshToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
}

func VerifyAccessToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
}


