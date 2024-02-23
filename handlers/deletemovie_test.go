package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"moviepin/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeleteMovieHandler(t *testing.T) {
	server := gin.New()
	mockService := &mock.ServiceMock{}

	handler := NewMovieHandler(mockService)

	route := "/movies/:movieID"
	routeHttpMethod := http.MethodDelete

	server.Handle(routeHttpMethod, route, handler.DeleteMovieHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		id     string
		err    mock.ErrMock
		status int
		resp   gin.H
	}{
		"request is successful": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusNoContent,
			resp:   gin.H{},
		},
		"not found": {
			id:     mock.Movie.ID,
			err:    mock.NotExistsError,
			status: http.StatusNoContent,
			resp:   gin.H{},
		},
		"db error": {
			id:     mock.Movie.ID,
			err:    mock.DeleteMovieError,
			status: http.StatusInternalServerError,
			resp: gin.H{
				"message": "failed to delete movie",
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

			if status := res.StatusCode; status != v.status {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)
			}

			// No need to check the body if the status code is 204
			if v.status == http.StatusNoContent {
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