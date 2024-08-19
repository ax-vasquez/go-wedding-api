package controllers

import (
	docs "github.com/ax-vasquez/wedding-site-api/docs"
	"github.com/ax-vasquez/wedding-site-api/middleware"
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

	// Routes without auth middleware (these are used to set/update the user's token, used by the auth middleware)
	{
		v1.POST("/signup", Signup)
		v1.POST("/login", Login)
	}

	// Routes for obtaining full or partial data sets for the base data types (admin-only)
	resourceRoutesV1 := v1.Group("")
	{
		resourceRoutesV1.Use(middleware.AuthenticateV1())
		resourceRoutesV1.GET("/entrees", middleware.IsAdminOrCurrentUser(), GetEntrees)
		resourceRoutesV1.GET("/users", middleware.IsAdmin(), GetUsers)
		resourceRoutesV1.GET("/horsdoeuvres", middleware.IsAdminOrCurrentUser(), GetHorsDoeuvres)
	}

	horsDoeuvresRoutesV1 := v1.Group("/horsdoeuvres")
	{
		horsDoeuvresRoutesV1.Use(middleware.AuthenticateV1())
		horsDoeuvresRoutesV1.POST("", middleware.IsAdmin(), CreateHorsDoeuvres)
		horsDoeuvresRoutesV1.DELETE("/:id", middleware.IsAdmin(), DeleteHorsDoeuvres)
	}

	entreeRoutesV1 := v1.Group("/entree")
	{
		entreeRoutesV1.Use(middleware.AuthenticateV1())
		entreeRoutesV1.POST("", middleware.IsAdmin(), CreateEntree)
		entreeRoutesV1.DELETE("/:id", middleware.IsAdmin(), DeleteEntree)
	}

	userRoutesV1 := v1.Group("/user")
	{
		userRoutesV1.Use(middleware.AuthenticateV1())
		userRoutesV1.GET("/:id/invitees", middleware.IsAdminOrCurrentUser(), GetInviteesForUser)
		userRoutesV1.GET("/:id/entrees", middleware.IsAdminOrCurrentUser(), GetEntrees)
		userRoutesV1.GET("/:id/horsdoeuvres", middleware.IsAdminOrCurrentUser(), GetHorsDoeuvres)
		userRoutesV1.PATCH("", middleware.IsAdminOrCurrentUser(), UpdateUser)
		userRoutesV1.POST("", middleware.IsAdmin(), CreateUser)
		userRoutesV1.POST("/:id/invite-user", middleware.IsAdminOrCurrentUser(), CreateUserInvitee)
		userRoutesV1.DELETE("/:id", middleware.IsAdmin(), DeleteUser)
	}

	inviteeRoutesV1 := v1.Group("/invitee")
	{
		inviteeRoutesV1.Use(middleware.AuthenticateV1())
		inviteeRoutesV1.DELETE("/:id", middleware.IsAdminOrCurrentUser(), DeleteInvitee)
	}

	return r
}

func SetupRoutes() error {
	r := paveRoutes()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run()
}
