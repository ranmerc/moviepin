package grpcclient

import (
	"context"
	"token-management-service/tokenpb"
)

type TokenServiceGRPCClient interface {
	GenerateToken(username string) (string, error)
	VerifyToken(token string) (bool, error)
}

type TokenServiceClient struct {
	client tokenpb.TokenServiceClient
}

func NewTokenServiceClient(client tokenpb.TokenServiceClient) *TokenServiceClient {
	return &TokenServiceClient{
		client: client,
	}
}

func (t *TokenServiceClient) GenerateToken(username string) (string, error) {
	req := &tokenpb.GenerateTokenRequest{
		Username: username,
	}

	res, err := t.client.GenerateToken(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func (t *TokenServiceClient) VerifyToken(token string) (bool, error) {
	req := &tokenpb.VerifyTokenRequest{
		Token: token,
	}

	res, err := t.client.VerifyToken(context.Background(), req)
	if err != nil {
		return false, err
	}

	return res.Valid, nil
}
