package schema

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	yamlData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(yamlData))

	yamlidatorURL := os.Getenv("YAMLIDATOR_URL")

	client := &http.Client{}

	response, _ := client.PostForm(yamlidatorURL+"/api/v1/schema", url.Values{"data": {sEnc}})

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Body.Close()

	ctx.Data(http.StatusOK, "text/html", body)

}
