package main

import (
	"featured-base-starter-kit/internal/config"
	"featured-base-starter-kit/internal/router"
	"featured-base-starter-kit/pkg/validation"
)

func main(){
	cfg := config.LoadConfig()

	validation.Init(cfg.Lang)

	r := router.Setup(&cfg)

	r.Run(":" + cfg.Port)
}