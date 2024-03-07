package grpcclient

import (
	"movie-management-service/mock"
	"testing"
)

func TestTokenServiceClient_GenerateToken(t *testing.T) {
	username := "username"

	cases := map[string]struct {
		mockErr mock.ErrMock
		token   string
		err     error
	}{
		"rpc call is successful": {
			mockErr: mock.OK,
			token:   "token",
			err:     nil,
		},
		"rpc call fails": {
			mockErr: mock.GenerateTokenError,
			token:   "",
			err:     mock.ErrMockTokenError,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			serviceClient := mock.NewTokenClientMock(tc.mockErr)

			client := NewTokenServiceClient(serviceClient)
			token, err := client.GenerateToken(username)

			if got, want := err, tc.err; got != want {
				t.Errorf("got error %v, want %v", got, want)
			}

			if got, want := token, tc.token; got != want {
				t.Errorf("got token %v, want %v", got, want)
			}

		})
	}
}

func TestTokenServiceClient_VerifyToken(t *testing.T) {
	cases := map[string]struct {
		mockErr mock.ErrMock
		valid   bool
		err     error
	}{
		"rpc call is successful": {
			mockErr: mock.OK,
			valid:   true,
			err:     nil,
		},
		"rpc call fails": {
			mockErr: mock.VerifyTokenError,
			valid:   false,
			err:     mock.ErrMockTokenError,
		},
		"rpc call returns invalid token": {
			mockErr: mock.VerifyTokenInvalidError,
			valid:   false,
			err:     nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			serviceClient := mock.NewTokenClientMock(tc.mockErr)

			client := NewTokenServiceClient(serviceClient)
			valid, err := client.VerifyToken("token")

			if got, want := err, tc.err; got != want {
				t.Errorf("got error %v, want %v", got, want)
			}

			if got, want := valid, tc.valid; got != want {
				t.Errorf("got valid %v, want %v", got, want)
			}

		})
	}
}
