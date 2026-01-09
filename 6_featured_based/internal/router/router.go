package router

import (
	"featured-base-starter-kit/internal/config"
	"featured-base-starter-kit/internal/middleware"
	"featured-base-starter-kit/internal/modules/user"

	"github.com/gin-gonic/gin"
)

func Setup (cfg *config.Config) *gin.Engine {
	//r := gin.Default()
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	//r.Use(gin.Logger())

	// api/users
	protectedRoute := r.Group("api")
	protectedRoute.Use(middleware.ApiKeyAuthMiddleware(cfg.API_SECRET_KEY))
	protectedRoute.POST("/users", user.CreateUserHandler)

	//r.POST("/users", user.CreateUserHandler)

	return r
}