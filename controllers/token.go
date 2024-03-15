package controllers

import (
	"time"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)


func createToken(username string, pmjay string, force_type string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &TokenClaim{
		Username: username,
		PmjayId: pmjay,
		ForceType: force_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (*TokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func getClaimsFromRequest(r *http.Request) (*TokenClaim, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	claims, err := validateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}