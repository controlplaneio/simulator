package router

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"wakeward/yaml-ctf/_app/pod-checker/frontend/app/about"
	"wakeward/yaml-ctf/_app/pod-checker/frontend/app/example"
	"wakeward/yaml-ctf/_app/pod-checker/frontend/app/home"
	"wakeward/yaml-ctf/_app/pod-checker/frontend/app/schema"
)

func New() *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/public", "frontend/static")
	router.LoadHTMLGlob("frontend/template/*")

	router.GET("/", home.Handler)
	router.GET("/about", about.Handler)
	router.POST("/schema", schema.Handler)
	router.POST("/example", example.Handler)

	return router
}
