package controllers

import "github.com/gin-gonic/gin"

type V1_API_RESPONSE struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}
