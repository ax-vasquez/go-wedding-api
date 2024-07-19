package routes

import (
	"github.com/ax-vasquez/wedding-site-api/controllers"
	"github.com/gin-gonic/gin"
)

func paveRoutes() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.DELETE("/user/:id", controllers.DeleteUser)
		v1.GET("/ping", controllers.Ping)
		v1.GET("/users", controllers.GetUsers)
		v1.PATCH("/user", controllers.UpdateUser)
		v1.POST("/user", controllers.CreateUser)
		v1.POST("/user/:id/invite-user", controllers.CreateUserInvitee)
	}

	return r
}

func Setup() error {
	r := paveRoutes()
	return r.Run()
}
