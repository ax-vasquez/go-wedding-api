package controllers

import (
	docs "github.com/ax-vasquez/wedding-site-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

type V1_API_RESPONSE struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

type DeleteRecordResponse struct {
	DeletedRecords int `json:"deleted_records"`
}

type V1_API_DELETE_RESPONSE struct {
	V1_API_RESPONSE
	Data DeleteRecordResponse `json:"data"`
}

func paveRoutes() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.DELETE("/entree/:id", DeleteEntree)
		v1.DELETE("/horsdoeuvres/:id", DeleteHorsDoeuvres)
		v1.DELETE("/user/:id", DeleteUser)
		v1.DELETE("/invitee/:invitee_id", DeleteInviteeForUser)
		v1.GET("/entrees", GetEntrees)
		v1.GET("/horsdoeuvres", GetHorsDoeuvres)
		v1.GET("/users", GetUsers)
		v1.GET("/user/:id/invitees", GetInviteesForUser)
		v1.GET("/user/:id/entrees", GetEntrees)
		v1.GET("/user/:id/horsdoeuvres", GetHorsDoeuvres)
		v1.PATCH("/user", UpdateUser)
		v1.POST("/entree", CreateEntree)
		v1.POST("/horsdoeuvres", CreateHorsDoeuvres)
		v1.POST("/user", CreateUser)
		v1.POST("/user/:id/invite-user", CreateUserInvitee)
	}

	return r
}

func SetupRoutes() error {
	r := paveRoutes()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run()
}
