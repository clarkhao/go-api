package main

import (
	"log"
	"os"

	"github.com/clarkhao/ws/router"
	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("PORT", "8080")
}

func main() {
	log.Printf("Gin cold start")
	r := gin.Default()
	// router group
	v0 := r.Group("/v0")
	{
		router.Routes(v0.Group(""))
	}
	log.Printf("start server on localhost:8080")
	//configuration
	gin.SetMode(gin.DebugMode)
	r.Run()
}
