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

func TestPutMovieHandler(t *testing.T) {
	server := gin.New()

	mockService := &mock.ServiceMock{}
	handler := NewMovieHandler(mockService)

	route := "/movies/:movieID"
	routeHttpMethod := http.MethodPut

	server.Handle(routeHttpMethod, route, handler.PutMovieHandler)
	httpServer := httptest.NewServer(server)

	body := gin.H{
		"ID":          mock.Movie.ID,
		"title":       mock.Movie.Title,
		"releaseDate": mock.Movie.ReleaseDate.Format(time.RFC3339),
		"genre":       mock.Movie.Genre,
		"director":    mock.Movie.Director,
		"description": mock.Movie.Description,
	}

	cases := map[string]struct {
		id     string
		err    mock.ErrMock
		status int
		body   gin.H
		resp   gin.H
	}{
		"movie put request is successful": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusOK,
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
		"movie put request not found when movie id is non-existent": {
			id:     mock.Movie.ID,
			err:    mock.NotExistsError,
			status: http.StatusNotFound,
			body:   body,
			resp: gin.H{
				"message": "movie does not exist",
			},
		},
		"movie put request when there is db error": {
			id:     mock.Movie.ID,
			err:    mock.UpdateMovieError,
			status: http.StatusInternalServerError,
			body:   body,
			resp: gin.H{
				"message": "failed to update movie",
			},
		},
		"movie put request when movie movie id is invalid": {
			id:     "invalid",
			err:    mock.OK,
			status: http.StatusBadRequest,
			body:   body,
			resp: gin.H{
				"message": "invalid id",
			},
		},
		"movie put request when date in body is empty": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusBadRequest,
			body:   gin.H{},
			resp: gin.H{
				"message": "invalid request body",
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
