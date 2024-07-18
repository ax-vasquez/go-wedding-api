package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	response := V1_API_RESPONSE{
		Status:  http.StatusOK,
		Message: "OK",
	}
	c.JSON(http.StatusOK, response)
}
