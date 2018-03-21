package main

import (
	"gopkg.in/hlandau/passlib.v1"
    "github.com/dgrijalva/jwt-go"
	"time"
)

func loginUser(login string, pass string) (*User, bool) {
	user := &User{
		User: login,
	}
	err := db.First(user).Error
	if err != nil {
		return nil, false
	}

	newHash, err := passlib.Verify(pass, user.Pass)
	if err != nil {
		return nil, false
	}

	if newHash != "" {
		user.Pass = newHash
		db.Save(user)
	}

	return user, true
}

func newJWT(user *User) (string, error) {
	claims := JWTClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour*24*7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString(JWTSecretBytes)
	if err != nil {
		return "", err
	}
	return s, nil
}

func refreshJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecretBytes, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return newJWT(claims.User)
	} else {
		return "", err
	}
}

func verifyJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecretBytes, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, nil
	}
}