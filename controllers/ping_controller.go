package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type V1_API_RESPONSE struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

func Ping(c *gin.Context) {
	response := V1_API_RESPONSE{
		Status:  http.StatusOK,
		Message: "OK",
	}
	c.JSON(http.StatusOK, response)
}
