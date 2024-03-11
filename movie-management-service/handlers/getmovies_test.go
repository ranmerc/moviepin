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
	"time"

	"github.com/gin-gonic/gin"
)

func TestGetMoviesHandler(t *testing.T) {
	server := gin.New()
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService, tokenClient)

	route := "/movies"
	routeHttpMethod := http.MethodGet

	server.Handle(routeHttpMethod, route, handler.GetMoviesHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		err    mock.ErrMock
		status int
		resp   gin.H
	}{
		"movies get request is successful": {
			err:    mock.OK,
			status: http.StatusOK,
			resp: gin.H{
				"movies": []gin.H{
					{
						"ID":          mock.Movie.ID,
						"title":       mock.Movie.Title,
						"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
						"genre":       mock.Movie.Genre,
						"director":    mock.Movie.Director,
						"description": mock.Movie.Description,
					},
				},
			},
		},
		"movies get request failed when there is db error": {
			err:    mock.GetMoviesError,
			status: http.StatusInternalServerError,
			resp: gin.H{
				"message": "failed to get movies",
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			if v.err != mock.OK {
				mockService.Err = v.err
			} else {
				mockService.Err = mock.OK
			}

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

			if status := res.StatusCode; status != v.status {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)
			}

			if fmt.Sprint(v.resp) != fmt.Sprint(got) {
				t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)
			}
		})
	}
}
