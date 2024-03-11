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

func TestDeleteMovieHandler(t *testing.T) {
	server := gin.New()
	mockService := &mock.ServiceMock{}
	mockClient := mock.NewTokenClientMock(mock.OK)
	tokenClient := grpcclient.NewTokenServiceClient(mockClient)

	handler := NewMovieHandler(mockService, tokenClient)

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
		"movie delete request is successful": {
			id:     mock.Movie.ID,
			err:    mock.OK,
			status: http.StatusNoContent,
			resp:   gin.H{},
		},
		"movie delete request successful when movie id is non-existent": {
			id:     mock.Movie.ID,
			err:    mock.DeleteMovieNotExistsError,
			status: http.StatusNoContent,
			resp:   gin.H{},
		},
		"movie delete request failed when there is db error": {
			id:     mock.Movie.ID,
			err:    mock.DeleteMovieError,
			status: http.StatusInternalServerError,
			resp: gin.H{
				"message": "failed to delete movie",
			},
		},
		"movie delete request failed when movie id is invalid": {
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
