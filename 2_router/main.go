package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users", func(ctx *gin.Context) {
		// active 
		// role
		// Query Parametresi
		// /users?active=true
		// /users?active=true&role=admin

		active := ctx.Query("active")
		role := ctx.Query("role")

		ctx.JSON(http.StatusOK, gin.H{
			"endpoint": "/users",
			"method" : "GET",
			"active": active,
			"role": role,
			"message": "Kullanıcı listesi (query param örneği)",
		})
	})

	router.GET("/users/:id", func(ctx *gin.Context) {
		// PATH parametre
		// /users/20

		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"endpoint": "/users/:id",
			"method": "GET",
			"user_id": id,
			"message": "Tek kullanıcı (path param örneği)",
		})
	})

	router.GET("/users/:id/details", func(ctx *gin.Context) {
		id := ctx.Param("id")
		isActive := ctx.Query("is_active")

		ctx.JSON(http.StatusOK, gin.H{
			"user_id": id,
			"is_active": isActive,
			"message": "Path + Query birlikte kullanım",
		})
	})

	router.Run(":8080")
}
