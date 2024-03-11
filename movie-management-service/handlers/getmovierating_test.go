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

func TestGetMovieRatingHandler(t *testing.T) {
	server := gin.New()
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService, tokenClient)

	route := "/movies/:movieID/rating"
	routeHttpMethod := http.MethodGet

	server.Handle(routeHttpMethod, route, handler.GetMovieRatingHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		id     string
		err    mock.ErrMock
		status int
		resp   gin.H
	}{
		"movie rating get request is successful": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusOK,
			resp: gin.H{
				"ID":          mock.MovieReview.ID,
				"title":       mock.MovieReview.Title,
				"releaseDate": mock.MovieReview.ReleaseDate.Format(time.RFC3339),
				"genre":       mock.MovieReview.Genre,
				"director":    mock.MovieReview.Director,
				"description": mock.MovieReview.Description,
				"rating":      mock.MovieReview.Rating,
			},
		},
		"movie rating get request not found when movie id is non-existent": {
			id:     mock.Movie.ID,
			err:    mock.GetMovieRatingNotExistsError,
			status: http.StatusNotFound,
			resp: gin.H{
				"message": "movie does not exist",
			},
		},
		"movie rating get request when there is db error": {
			id:     mock.Movie.ID,
			err:    mock.GetMovieRatingError,
			status: http.StatusInternalServerError,
			resp: gin.H{
				"message": "failed to get movie rating",
			},
		},
		"movie rating get request when movie movie id is invalid": {
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
			requestURL := httpServer.URL + fmt.Sprintf("/movies/%s/rating", v.id)
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
