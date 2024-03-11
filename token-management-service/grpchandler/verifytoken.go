package grpchandler

import (
	"context"
	"movie-management-service/utils"
	"token-management-service/tokenpb"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken verifies the token
func (h Handler) VerifyToken(_ context.Context, req *tokenpb.VerifyTokenRequest) (*tokenpb.VerifyTokenResponse, error) {
	token := req.GetToken()

	if len(token) == 0 {
		return nil, ErrEmptyToken
	}

	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}
	withValidMethods := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})

	if _, err := jwt.ParseWithClaims(token, claims, keyFunc, withValidMethods); err != nil {
		utils.ErrorLogger.Print(err)
		return &tokenpb.VerifyTokenResponse{
			Valid:    false,
			Username: "",
		}, nil
	}

	username, err := claims.GetSubject()
	if err != nil {
		utils.ErrorLogger.Print(err)
		return nil, ErrInternal
	}

	return &tokenpb.VerifyTokenResponse{
		Valid:    true,
		Username: username,
	}, nil
}
