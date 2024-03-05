package grpchandler

import (
	"context"
	"reflect"
	"testing"
	"token-management-service/tokenpb"
)

func TestVerifyToken(t *testing.T) {
	h := Handler{}
	username := "username"

	res, err := h.GenerateToken(context.Background(), &tokenpb.GenerateTokenRequest{
		Username: username,
	})
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		token string
		res   *tokenpb.VerifyTokenResponse
		err   error
	}{
		"token is empty": {
			token: "",
			res:   nil,
			err:   ErrEmptyToken,
		},
		"token is invalid": {
			token: "token",
			res: &tokenpb.VerifyTokenResponse{
				Valid:    false,
				Username: "",
			},
			err: nil,
		},
		"token is valid": {
			token: res.Token,
			res: &tokenpb.VerifyTokenResponse{
				Valid:    true,
				Username: username,
			},
			err: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &tokenpb.VerifyTokenRequest{
				Token: tc.token,
			}

			h := Handler{}
			res, err := h.VerifyToken(context.Background(), req)

			if got, want := err, tc.err; got != want {
				t.Errorf("got %v, want %v", got, want)
				return
			}

			if got, want := res, tc.res; !reflect.DeepEqual(got, want) {
				t.Errorf("got %+v, want %+v", got, want)
			}
		})
	}
}
