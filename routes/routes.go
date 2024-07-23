package routes

import (
	"github.com/ax-vasquez/wedding-site-api/controllers"
	"github.com/gin-gonic/gin"
)

func paveRoutes() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.DELETE("/entree/:id", controllers.DeleteEntree)
		v1.DELETE("/horsdoeuvres/:id", controllers.DeleteHorsDoeuvres)
		v1.DELETE("/user/:id", controllers.DeleteUser)
		v1.DELETE("/user/:id/invitee/:invitee_id", controllers.DeleteInviteeForUser)
		v1.GET("/entrees", controllers.GetEntrees)
		v1.GET("/horsdoeuvres", controllers.GetHorsDoeuvres)
		v1.GET("/users", controllers.GetUsers)
		v1.GET("/user/:id/invitees", controllers.GetInviteesForUser)
		v1.GET("/user/:id/entrees", controllers.GetEntrees)
		v1.GET("/user/:id/horsdoeuvres", controllers.GetHorsDoeuvres)
		v1.PATCH("/user", controllers.UpdateUser)
		v1.POST("/entree", controllers.CreateEntree)
		v1.POST("/horsdoeuvres", controllers.CreateHorsDoeuvres)
		v1.POST("/user", controllers.CreateUsers)
		v1.POST("/user/:id/invite-user", controllers.CreateUserInvitee)
	}

	return r
}

func Setup() error {
	r := paveRoutes()
	return r.Run()
}
