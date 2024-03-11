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

func TestGetMovieHandler(t *testing.T) {
	server := gin.New()
	mockService := &mock.ServiceMock{}
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	handler := NewMovieHandler(mockService, tokenClient)

	route := "/movies/:movieID"
	routeHttpMethod := http.MethodGet

	server.Handle(routeHttpMethod, route, handler.GetMovieHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		id     string
		err    mock.ErrMock
		status int
		resp   gin.H
	}{
		"movie get request is successful": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusOK,
			resp: gin.H{
				"movie": gin.H{
					"ID":          mock.Movie.ID,
					"title":       mock.Movie.Title,
					"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
					"genre":       mock.Movie.Genre,
					"director":    mock.Movie.Director,
					"description": mock.Movie.Description,
				},
			},
		},
		"movie get request failed when movie id is non existent": {
			id:     mock.Movie.ID,
			err:    mock.GetMovieNotExistsError,
			status: http.StatusNotFound,
			resp: gin.H{
				"message": "movie does not exist",
			},
		},
		"movie get request failed when there is db error": {
			id:     mock.Movie.ID,
			err:    mock.GetMovieError,
			status: http.StatusInternalServerError,
			resp: gin.H{
				"message": "failed to get movie",
			},
		},
		"movie get request failed when movie id is invalid": {
			id:     "invalid",
			err:    mock.OK,
			status: http.StatusBadRequest,
			resp: gin.H{
				"message": gin.H{
					"movieID": "should be an UUID",
				},
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
			requestURL := httpServer.URL + fmt.Sprintf("/movies/%s", v.id)
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
