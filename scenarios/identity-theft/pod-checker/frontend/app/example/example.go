package example

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	yamlidatorURL := os.Getenv("YAMLIDATOR_URL")

	client := &http.Client{}

	response, _ := client.Post(yamlidatorURL+"/api/v1/example", "application/json", bytes.NewBuffer(jsonData))

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Body.Close()

	ctx.Data(http.StatusOK, "text/html", body)

}
