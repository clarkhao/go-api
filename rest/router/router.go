package router

import (
	"time"

	"restapi/controller/ping"
	"restapi/controller/stream"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

// Routes handle router in web app
func Routes(g *gin.RouterGroup) {
	// caching
	store := persistence.NewInMemoryStore(time.Minute)
	// ping router
	g.GET("/ping", logger.SetLogger(), cache.CachePage(store, time.Minute, ping.GetHandler))
	g.POST("/ping", logger.SetLogger(), cache.CachePage(store, time.Minute, ping.PostHandler))
	g.GET("/ping/err", logger.SetLogger(), ping.ErrorHandler)
	// with client
	g.GET("/ping/getclient", logger.SetLogger(), ping.ClientGetHandler)
	g.POST("/ping/postclient", logger.SetLogger(), ping.ClientPostHandler)
	// streaming
	g.GET("/stream", logger.SetLogger(), stream.GetChunkHandler)
	g.GET("/stream/:filename", logger.SetLogger(), stream.GetStreamHandler)
}
