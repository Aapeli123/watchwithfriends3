package server

import "github.com/gin-gonic/gin"

func fuckOffMyServer(g *gin.Context) {
	g.String(404, "Can you please fuck off my server")
}

func defineRoutes(r *gin.Engine) {
	r.GET("/isOnline", alive)
	r.GET("/ws", handleWS)
	r.NoRoute(fuckOffMyServer)
}
