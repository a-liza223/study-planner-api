package handlers

import (
	"net/http"
	"study-tracker-api/storage"

	"github.com/gin-gonic/gin"
)

func GetWorkload(c *gin.Context) {
	workload := storage.GetWorkload()
	c.JSON(http.StatusOK, workload)
}
