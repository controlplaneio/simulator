package about

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for our home page.
func Handler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "about.html", nil)
}
