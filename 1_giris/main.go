package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	/*router2 := gin.New()
	router2.Use(gin.Logger())
	router2.Use(gin.Recovery())*/

	router.GET("/", func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{
			"message": "Merhaba gin framework",
			"message2": 22,
		})
		// Content:Type: application/json
		// map[string]string
		// map[string]interface{}
	})

	router.GET("/panic", func (c *gin.Context)  {
		panic("Bilerek painic attık (Recovery demo)")
	})

	//port := ":8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port) // Network address tanımlaması yapıyorum
	// host:port

	if err := router.Run(addr); err != nil {
		panic(err)
	}
}
