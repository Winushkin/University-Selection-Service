package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestServer_generateAccessToken(t *testing.T) {
	server := &Server{
		jwtSecret: "testsecret",
	}
	userID := 1
	tokenStr, err := server.generateAccessToken(userID)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}
	if tokenStr == "" {
		t.Error("expected non-empty token")
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(server.jwtSecret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !token.Valid {
		t.Error("token is not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("failed to get claims")
	}
	if claims["user_id"] != float64(userID) {
		t.Errorf("expected user_id %d, got %v", userID, claims["user_id"])
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("exp claim is not a number")
	}
	if exp <= float64(time.Now().Unix()) {
		t.Error("exp is not in the future")
	}
}

func TestServer_generateRefreshToken(t *testing.T) {
	server := &Server{}
	tokenStr, err := server.generateRefreshToken()
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
	if tokenStr == "" {
		t.Error("expected non-empty token")
	}
	if len(tokenStr) != 44 {
		t.Errorf("expected token length 44, got %d", len(tokenStr))
	}
}
