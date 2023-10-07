package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Kind struct {
	ID   string `json:"id,omitempty"`
	Kind string `json:"kind"`
}

func main() {
	router := gin.Default()
	router.GET("/", run)
	router.Run("localhost:8080")
}

func run(c *gin.Context) {
	var flag = os.Getenv("FLAG")
	c.JSON(http.StatusOK, "Tokungfu Shop Running: "+flag)
}
