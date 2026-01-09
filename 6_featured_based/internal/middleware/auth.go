package middleware

import (
	"featured-base-starter-kit/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuthMiddleware(validSecretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")

		if apiKey != validSecretKey {
			// Abort request zincirini iptal
			api.SendError(ctx, http.StatusUnauthorized, "Yetkisiz erişim. Lütfen bilgileri kontrol edin", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}