package interceptors

import (
	"context"
	"testing"
	"time"

	"University-Selection-Service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func mockHandler(_ context.Context, _ interface{}) (interface{}, error) {
	return "response", nil
}

func TestAuthInterceptor_PublicMethods(t *testing.T) {
	cfg := &config.UserConfig{JWTSecret: "testsecret"}
	interceptor := AuthInterceptor(cfg)

	tests := []struct {
		method string
	}{
		{"/api.UserService/SignUp"},
		{"/api.UserService/Refresh"},
		{"/api.UserService/Logout"},
		{"/api.UserService/Login"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			ctx := context.Background()
			info := &grpc.UnaryServerInfo{FullMethod: tt.method}
			resp, err := interceptor(ctx, nil, info, mockHandler)
			assert.NoError(t, err)
			assert.Equal(t, "response", resp)
		})
	}
}

func TestAuthInterceptor_ProtectedMethods(t *testing.T) {
	cfg := &config.UserConfig{JWTSecret: "testsecret"}
	interceptor := AuthInterceptor(cfg)
	method := "/api.UserService/ProtectedMethod"

	t.Run("missing metadata", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: method}
		_, err := interceptor(ctx, nil, info, mockHandler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Contains(t, err.Error(), "Missing metadata")
	})

	t.Run("missing authorization header", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
		info := &grpc.UnaryServerInfo{FullMethod: method}
		_, err := interceptor(ctx, nil, info, mockHandler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Contains(t, err.Error(), "Missing Authorization header")
	})

	t.Run("invalid authorization format", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "InvalidToken"))
		info := &grpc.UnaryServerInfo{FullMethod: method}
		_, err := interceptor(ctx, nil, info, mockHandler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Contains(t, err.Error(), "Invalid Authorization header: must start with Bearer")
	})

	t.Run("invalid token", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer invalidtoken"))
		info := &grpc.UnaryServerInfo{FullMethod: method}
		_, err := interceptor(ctx, nil, info, mockHandler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("valid token", func(t *testing.T) {
		token, _ := generateTestToken(cfg, 1)
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		info := &grpc.UnaryServerInfo{FullMethod: method}
		resp, err := interceptor(ctx, nil, info, func(ctx context.Context, req interface{}) (interface{}, error) {
			userID, ok := ctx.Value("user_id").(int)
			assert.True(t, ok)
			assert.Equal(t, 1, userID)
			return "response", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})
}

func TestValidateToken(t *testing.T) {
	cfg := &config.UserConfig{JWTSecret: "testsecret"}

	t.Run("valid token", func(t *testing.T) {
		token, _ := generateTestToken(cfg, 1)
		id, err := validateToken(token, cfg)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("expired token", func(t *testing.T) {
		token, _ := generateTestTokenWithExp(cfg, 1, time.Now().Add(-1*time.Hour))
		_, err := validateToken(token, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validateToken: token has invalid claims: token is expired")
	})

	t.Run("invalid signing method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
		})
		tokenStr, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		_, err := validateToken(tokenStr, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValidateToken: unexpected signing method")
	})

	t.Run("invalid user id", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "invalid",
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
		})
		tokenStr, _ := token.SignedString([]byte(cfg.JWTSecret))
		_, err := validateToken(tokenStr, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validateToken: invalid user id")
	})
}

// Helper functions for generating test tokens
func generateTestToken(cfg *config.UserConfig, userID int) (string, error) {
	return generateTestTokenWithExp(cfg, userID, time.Now().Add(15*time.Minute))
}

func generateTestTokenWithExp(cfg *config.UserConfig, userID int, exp time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
	})
	return token.SignedString([]byte(cfg.JWTSecret))
}
