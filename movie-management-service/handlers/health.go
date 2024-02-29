package handlers

import (
	"movie-management-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mh MovieHandler) HealthCheckHandler(c *gin.Context) {
	if _, err := mh.domain.DBStatus(); err != nil {
		c.JSON(http.StatusFailedDependency, &model.HeathResponse{
			Status:   "alive",
			DBStatus: false,
		})
		return
	}

	c.JSON(http.StatusOK, &model.HeathResponse{
		Status:   "alive",
		DBStatus: true,
	})
}
