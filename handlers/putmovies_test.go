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

func TestPutMoviesHandler(t *testing.T) {
	server := gin.New()

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService)

	route := "/movies"
	routeHttpMethod := http.MethodPut

	server.Handle(routeHttpMethod, route, handler.PutMoviesHandler)
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
		"movies put request is successful": {
			err:    mock.OK,
			status: http.StatusNoContent,
			body:   body,
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
		"movies put request failed when there is db error": {
			err:    mock.ReplaceMoviesError,
			status: http.StatusInternalServerError,
			body:   body,
			resp: gin.H{
				"message": "failed to replace movies",
			},
		},
		"movies put request failed when request body is empty": {
			err:    mock.OK,
			status: http.StatusBadRequest,
			body:   gin.H{},
			resp: gin.H{
				"message": "invalid request body",
			},
		},
		"movies put request when there are no movies in request": {
			err:    mock.OK,
			status: http.StatusBadRequest,
			body:   gin.H{"movies": []gin.H{}},
			resp: gin.H{
				"message": "at least one movie is required",
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

			if status := res.StatusCode; status != v.status {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)
			}

			// If status is 204 or 404, there is no response body.
			if v.status == http.StatusNoContent || v.status == http.StatusNotFound {
				return
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

			if fmt.Sprint(v.resp) != fmt.Sprint(got) {
				t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)
			}
		})
	}
}
