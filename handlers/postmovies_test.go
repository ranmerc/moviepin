package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"moviepin/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestPostMoviesHandler(t *testing.T) {
	server := gin.New()

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService)

	route := "/movies"
	routeHttpMethod := http.MethodPost

	server.Handle(routeHttpMethod, route, handler.PostMoviesHandler)
	httpServer := httptest.NewServer(server)

	body := gin.H{
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
	}

	cases := map[string]struct {
		err    mock.ErrMock
		status int
		body   gin.H
		resp   gin.H
	}{
		"all movies added": {
			err:    mock.OK,
			status: http.StatusCreated,
			body:   body,
			resp: gin.H{
				"addedMovies": []gin.H{
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
		"all movies failed": {
			err:    mock.AddMovieError,
			status: http.StatusInternalServerError,
			body:   body,
			resp: gin.H{
				"message": "failed to add movies",
			},
		},
		"some movies failed": {
			err:    mock.OK,
			status: http.StatusMultiStatus,
			body: gin.H{
				"movies": []gin.H{
					{
						"ID":          mock.Movie.ID,
						"title":       mock.Movie.Title,
						"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
						"genre":       mock.Movie.Genre,
						"director":    mock.Movie.Director,
						"description": mock.Movie.Description,
					},
					{
						"ID":          mock.MovieIDToFail,
						"title":       mock.Movie.Title,
						"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
						"genre":       mock.Movie.Genre,
						"director":    mock.Movie.Director,
						"description": mock.Movie.Description,
					},
				},
			},
			resp: gin.H{
				"addedMovies": []gin.H{
					{
						"ID":          mock.Movie.ID,
						"title":       mock.Movie.Title,
						"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
						"genre":       mock.Movie.Genre,
						"director":    mock.Movie.Director,
						"description": mock.Movie.Description,
					},
				},
				"failedMovies": []gin.H{
					{
						"ID":          mock.MovieIDToFail,
						"title":       mock.Movie.Title,
						"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
						"genre":       mock.Movie.Genre,
						"director":    mock.Movie.Director,
						"description": mock.Movie.Description,
					},
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
			requestURL := httpServer.URL + route

			jsonBody, err := json.Marshal(v.body)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			reqBody := bytes.NewBuffer(jsonBody)

			req, err := http.NewRequest(routeHttpMethod, requestURL, reqBody)
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
			err = json.Unmarshal(body, &got)
			if err != nil {
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
