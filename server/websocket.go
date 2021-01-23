package server

import (
	"fmt"
	"net/http"

	"github.com/Aapeli123/watchwithfriends3/lib/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsConnection struct {
	Conn    *websocket.Conn
	Session *session.Session
}

type wsMessage struct {
	Operation string
	Data      interface{}
}

func handleMessage(msg wsMessage, wsConn wsConnection) {

}

func handleWS(c *gin.Context) {
	SID, err := getReqSID(c)
	sidHeader := http.Header{}

	valid := session.ValidateSess(SID)

	if err == http.ErrNoCookie || !valid {
		SID = session.AddSess().ID
		sidCookie := http.Cookie{Name: "SID", Value: SID, HttpOnly: true}
		sidHeader.Add("Set-Cookie", sidCookie.String())
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, sidHeader)

	if err != nil {
		fmt.Println("Connection upgrade failed...")
	}
	sess, err := session.GetSess(SID)
	wsConn := wsConnection{Conn: conn, Session: &sess}
	msg := wsMessage{}

	for {
		if wsConn.Conn.ReadJSON(&msg) != nil {
			break
		}
		handleMessage(msg, wsConn)
	}
}
