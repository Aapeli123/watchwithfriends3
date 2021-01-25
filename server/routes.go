package server

import "github.com/gin-gonic/gin"

func defineRoutes(r *gin.Engine) {
	r.GET("/isOnline", alive)
	r.GET("/ws", handleWS)
}
