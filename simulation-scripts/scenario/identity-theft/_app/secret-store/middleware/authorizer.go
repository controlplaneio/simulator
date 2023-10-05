package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"simulation-scripts/scenario/identity-theft/_app/secret-store/config"
)

// Authorize ID Token and Claims.
func Authorizer(auth *config.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			ctx.String(http.StatusForbidden, "No Authorization header provided")
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			ctx.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			ctx.Abort()
			return
		}
		idToken, err := auth.IDTokenVerifier.Verify(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusForbidden, "Bearer token is invalid")
			ctx.Abort()
			return
		}

		var claims struct {
			Email    string   `json:"email"`
			Verified bool     `json:"email_verified"`
			Groups   []string `json:"groups"`
		}
		if err := idToken.Claims(&claims); err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to parse claims")
			ctx.Abort()
			return
		}

		if !claims.Verified {
			ctx.String(http.StatusForbidden, "Email (%q) is not verified", claims.Email)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
