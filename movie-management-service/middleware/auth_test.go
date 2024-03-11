package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"movie-management-service/grpcclient"
	"movie-management-service/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuth_Middleware(t *testing.T) {
	router := gin.Default()
	mockClient := mock.NewTokenClientMock(mock.OK)

	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	auth := NewAuthMiddleware(tokenClient)

	route := "/test"
	routeHttpMethod := http.MethodGet
	handler := func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}

	router.Handle(routeHttpMethod, route, auth.Middleware(), handler)

	server := httptest.NewServer(router)

	defer server.Close()

	tests := []struct {
		name   string
		err    mock.ErrMock
		status int
		header string
		resp   gin.H
	}{
		{
			name:   "auth validation is successful",
			err:    mock.OK,
			status: http.StatusOK,
			header: "Bearer token",
			resp: gin.H{
				"message": "success",
			},
		},
		{
			name:   "auth validation failed when token is invalid",
			err:    mock.VerifyTokenInvalidError,
			status: http.StatusUnauthorized,
			header: "Bearer token",
			resp: gin.H{
				"message": ErrInvalidToken.Error(),
			},
		},
		{
			name:   "auth validation failed when auth header is not provided",
			err:    mock.OK,
			status: http.StatusUnauthorized,
			header: "",
			resp: gin.H{
				"message": ErrMissingAuthHeader.Error(),
			},
		},
		{
			name:   "auth validation failed when auth header is invalid",
			err:    mock.OK,
			status: http.StatusUnauthorized,
			header: "token",
			resp: gin.H{
				"message": ErrInvalidAuthHeader.Error(),
			},
		},
	}

	for _, v := range tests {
		mockClient.Err = v.err

		request, _ := http.NewRequest(routeHttpMethod, server.URL+route, nil)

		if len(v.header) != 0 {
			request.Header.Set("Authorization", v.header)
		}

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatalf("unable to make http request %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error("unexpected error: ", err)
		}

		var got gin.H
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != v.status {
			t.Fatalf("expected status code %d, got %d", v.status, res.StatusCode)
		}

		if fmt.Sprint(v.resp) != fmt.Sprint(got) {
			t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)
		}
	}
}
