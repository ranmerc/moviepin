package grpchandler

import (
	"context"
	"testing"
	"token-management-service/tokenpb"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	cases := map[string]struct {
		username string
		err      error
	}{
		"username is not empty": {
			username: "username",
			err:      nil,
		},
		"username is empty": {
			username: "",
			err:      ErrEmptyUsername,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &tokenpb.GenerateTokenRequest{
				Username: tc.username,
			}

			h := Handler{}
			res, err := h.GenerateToken(context.Background(), req)

			if tc.err != nil {
				if err != tc.err {
					t.Errorf("got error: %v, want: %v", err, tc.err)
				}
			} else {
				claims := &Claims{}
				keyFunc := func(token *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				}
				withValidMethod := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})

				if _, err := jwt.ParseWithClaims(res.Token, claims, keyFunc, withValidMethod); err != nil {
					t.Errorf("got invalid token: %v", err)
				}

				if claims.Subject != tc.username {
					t.Errorf("got username: %s, want: %s", claims.Subject, tc.username)
				}
			}
		})

	}
}
