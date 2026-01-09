package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		log.Printf("[GIRIS] %s isteği %s adresinden %s yoluna geldi.", method, clientIP, path)

		// Zincirdeki bir sonraki Middleware e ya da middleware yoksa ana fonksiyona yani handlera gider.
		ctx.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime) // işlem süresi işlem ne kadar sürdü

		statusCode := ctx.Writer.Status()

		log.Printf("[CIKIS] Durum: %d | Süre: %v\n", statusCode, latency)
	}
}