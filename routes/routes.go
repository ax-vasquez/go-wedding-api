package routes

import (
	"github.com/ax-vasquez/wedding-site-api/controllers"
	"github.com/gin-gonic/gin"
)

func paveRoutes() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", controllers.Ping)
	}

	return r
}

func Setup() error {
	r := paveRoutes()
	return r.Run()
}
