package server

import "github.com/go-gin/gin"

func defineRoutes(r *gin.Engine) {
	r.GET("/")
	r.GET("/ws", handleWS)
}
