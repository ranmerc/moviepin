package mock

import (
	"context"
	"token-management-service/tokenpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrMockTokenError is a mock error for token client.
	ErrMockTokenError = status.Errorf(codes.Unknown, "mock error")
)

const (
	// GenerateTokenError is a mock error for GenerateToken.
	GenerateTokenError ErrMock = iota

	// VerifyTokenError is a mock error for VerifyToken.
	VerifyTokenError

	// VerifyTokenInvalidError is a mock error for VerifyToken when token is invalid.
	VerifyTokenInvalidError
)

const (
	// token is a mock token.
	Token = "token"

	// username is a mock username.
	Username = "username"
)

type TokenClientMock struct {
	Err ErrMock
}

func NewTokenClientMock(err ErrMock) *TokenClientMock {
	return &TokenClientMock{
		Err: err,
	}
}

func (t TokenClientMock) GenerateToken(ctx context.Context, in *tokenpb.GenerateTokenRequest, opts ...grpc.CallOption) (*tokenpb.GenerateTokenResponse, error) {
	if t.Err == GenerateTokenError {
		return nil, ErrMockTokenError
	}

	return &tokenpb.GenerateTokenResponse{
		Token: Token,
	}, nil
}

func (t TokenClientMock) VerifyToken(ctx context.Context, in *tokenpb.VerifyTokenRequest, opts ...grpc.CallOption) (*tokenpb.VerifyTokenResponse, error) {
	if t.Err == VerifyTokenError {
		return nil, ErrMockTokenError
	}

	if t.Err == VerifyTokenInvalidError {
		return &tokenpb.VerifyTokenResponse{
			Valid:    false,
			Username: "",
		}, nil
	}

	return &tokenpb.VerifyTokenResponse{
		Valid:    true,
		Username: Username,
	}, nil
}
