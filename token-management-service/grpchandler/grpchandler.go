package grpchandler

import (
	"token-management-service/tokenpb"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// secretKey is the secret key for jwt token.
	secretKey = []byte("secret")

	// ErrEmptyUsername is gRPC error when in request username is empty
	ErrEmptyUsername = status.Error(codes.InvalidArgument, "username cannot be empty")

	// ErrEmptyToken is gRPC error when in request token is empty
	ErrEmptyToken = status.Error(codes.InvalidArgument, "token cannot be empty")

	// ErrInternal is gRPC error when something went wrong
	ErrInternal = status.Error(codes.Internal, "something went wrong")
)

// Claims is custom claims for jwt token
type Claims struct {
	jwt.RegisteredClaims
}

// Handler implements tokenpb.TokenServiceServer
type Handler struct {
	tokenpb.UnimplementedTokenServiceServer
}
