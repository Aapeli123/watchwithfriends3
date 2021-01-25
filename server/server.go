package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Options struct {
	Address string
	HTTPS   bool
	Cert    string
	Key     string
}

var upgrader websocket.Upgrader

func Start(options Options) {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	r := gin.Default()
	c := cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowAllOrigins:  false,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowHeaders:     []string{"content-type", "content-length"},
	})
	r.Use(c)

	defineRoutes(r)
	if options.HTTPS {
		r.RunTLS(options.Address, options.Cert, options.Key)
	} else {
		r.Run(options.Address)
	}
}

func getReqSID(c *gin.Context) (string, error) {
	sid, err := c.Request.Cookie("SID")
	if err != nil {
		return "", err
	}
	return sid.Value, nil
}

func alive(c *gin.Context) {
	c.String(200, "true")
}
