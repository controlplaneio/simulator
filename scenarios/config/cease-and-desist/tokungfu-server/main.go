package main

import (
	"log"
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

	log.Print("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe("0.0.0.0:8080", router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}

func run(c *gin.Context) {
	var flag = os.Getenv("FLAG")
	c.JSON(http.StatusOK, "Tokungfu Shop Running: "+flag)
}
