package handlers

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

func TestHealthCheckHandler(t *testing.T) {
	server := gin.New()
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService, tokenClient)

	route := "/health"
	routeHttpMethod := http.MethodGet

	server.Handle(routeHttpMethod, route, handler.HealthCheckHandler)

	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		dbErr mock.ErrMock
		code  int
		want  gin.H
	}{
		"health API responded with DB status": {
			dbErr: mock.OK,
			code:  http.StatusOK,
			want: gin.H{
				"status": "alive",
				"db":     true,
			},
		},
		"health API responded with DB error": {
			dbErr: mock.DBStatusError,
			code:  http.StatusFailedDependency,
			want: gin.H{
				"status": "alive",
				"db":     false,
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			mockService.Err = v.dbErr

			client := http.Client{}

			requestURL := httpServer.URL + route
			req, err := http.NewRequest(routeHttpMethod, requestURL, nil)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

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

			if status := res.StatusCode; status != v.code {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.code)
			}

			if fmt.Sprint(v.want) != fmt.Sprint(got) {
				t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.want)
			}
		})
	}
}
