package controllers

import (
	"wakeward/yaml-ctf/_app/secret-store/config"
	docs "wakeward/yaml-ctf/_app/secret-store/docs"
	"wakeward/yaml-ctf/_app/secret-store/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func UserRoutes(auth *config.Authenticator, r *gin.Engine) {

	docs.SwaggerInfo.BasePath = "api/v1"
	v1 := r.Group("/api/v1")
	{
		p := v1.Group("/")
		p.Use(middleware.Authorizer(auth))
		{
			p.GET("/users", GetUsers)
			p.GET("/users/:userId", GetUserByID)
			p.POST("/users", NewUser)
			p.PUT("/users/:userId", UpdateUser)
			p.DELETE("/users/:userId", DeleteUser)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
