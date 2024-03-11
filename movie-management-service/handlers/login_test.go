package handlers

import (
	"bytes"
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

func TestLoginHandler(t *testing.T) {
	server := gin.New()
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	mockService := mock.NewServiceMock()
	handler := NewMovieHandler(&mockService, tokenClient)

	route := "/login"
	routeHttpMethod := http.MethodPost

	server.Handle(routeHttpMethod, route, handler.LoginHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		err     mock.ErrMock
		grpcErr mock.ErrMock
		status  int
		req     string
		resp    gin.H
	}{
		"login request is successful": {
			err:     mock.OK,
			grpcErr: mock.OK,
			status:  http.StatusOK,
			req:     "username=username&password=password",
			resp: gin.H{
				"token": mock.Token,
			},
		},
		"login request failed when credentials are invalid": {
			err:     mock.InvalidCredentialsError,
			grpcErr: mock.OK,
			status:  http.StatusUnauthorized,
			req:     "username=username&password=password",
			resp: gin.H{
				"message": "invalid credentials",
			},
		},
		"login request failed when db error occurs": {
			err:     mock.LoginUserError,
			grpcErr: mock.OK,
			status:  http.StatusInternalServerError,
			req:     "username=username&password=password",
			resp: gin.H{
				"message": "failed to login",
			},
		},
		"login request failed when username is too short": {
			err:     mock.OK,
			grpcErr: mock.OK,
			status:  http.StatusBadRequest,
			req:     "username=user&password=password",
			resp: gin.H{
				"message": gin.H{
					"username": "should be minimum 6 characters",
				},
			},
		},
		"login request failed when password is too short": {
			err:     mock.OK,
			grpcErr: mock.OK,
			status:  http.StatusBadRequest,
			req:     "username=username&password=pass",
			resp: gin.H{
				"message": gin.H{
					"password": "should be minimum 8 characters",
				},
			},
		},
		"login request failed when unable to generate token": {
			err:     mock.OK,
			grpcErr: mock.GenerateTokenError,
			status:  http.StatusInternalServerError,
			req:     "username=username&password=password",
			resp: gin.H{
				"message": mock.ErrMockTokenError.Error(),
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			mockClient.Err = v.grpcErr
			mockService.Err = v.err

			client := http.Client{}
			requestURL := httpServer.URL + route

			reqBody := bytes.NewBufferString(v.req)
			req, err := http.NewRequest(routeHttpMethod, requestURL, reqBody)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res, err := client.Do(req)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			var got gin.H
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatal(err)
			}

			if status := res.StatusCode; status != v.status {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)
			}

			if fmt.Sprint(v.resp) != fmt.Sprint(got) {
				t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)
			}
		})
	}
}
