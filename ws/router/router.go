package router

import (
	"github.com/clarkhao/ws/controller/ping"
	"github.com/gin-gonic/gin"
)

// Routes handle router in web app
func Routes(g *gin.RouterGroup) {
	g.GET("/ws", ping.WSConnHandler)
}
