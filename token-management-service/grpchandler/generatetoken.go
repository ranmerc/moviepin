package grpchandler

import (
	"context"
	"movie-management-service/utils"
	"time"
	"token-management-service/tokenpb"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates jwt token with username
func (h Handler) GenerateToken(_ context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	username := req.GetUsername()

	if len(username) == 0 {
		return nil, ErrEmptyUsername
	}

	claims := &Claims{
		jwt.RegisteredClaims{
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.ErrorLogger.Print(err)
		return nil, ErrInternal
	}

	return &tokenpb.GenerateTokenResponse{
		Token: signedToken,
	}, nil
}
