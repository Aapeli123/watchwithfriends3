package server

import (
	"fmt"

	"github.com/Aapeli123/watchwithfriends3/lib/room"
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
	Data      msgData
}

type msgData struct {
	Username    string
	UserID      string
	Room        string
	VideoLink   string
	Time        float64
	ChatMessage string
}

type wsResponse struct {
	Operation string
	Success   bool
	Data      interface{}
}

func sync(rID string, wsConn *wsConnection, time float64) {
	r, err := room.GetRoom(wsConn.Session.InRoom)
	if wsConn.Session.IsLeader(rID) {
		if err == nil {
			r.UpdateVideoTime(time)
		}
	}

	d, err := getRoomData(r)
	if err != nil {
		return
	}
	wsConn.Conn.WriteJSON(wsResponse{
		Operation: "sync",
		Data:      d,
		Success:   true,
	})

}

func changeRoomVideo(rID string, wsConn *wsConnection, vidLink string) {
	r, err := room.GetRoom(wsConn.Session.InRoom)
	if err == nil {
		r.UpdateVideoLink(vidLink)
	}
}

type roomData struct {
	room.Room
	Usernames []string
}

func getRoomData(r room.Room) (roomData, error) {
	data := roomData{
		Room:      r,
		Usernames: []string{},
	}
	for _, u := range r.Users {
		s, err := session.GetSess(u)
		un := s.Username
		if err != nil {
			return roomData{}, err
		}
		data.Usernames = append(data.Usernames, un)
	}
	return data, nil
}

func joinRoom(rID string, wsConn *wsConnection) {
	wsConn.Session.JoinRoom(rID)
	r, err := room.GetRoom(rID)
	if err != nil {
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "join",
			Success:   false,
		})
		return
	}
	d, err := getRoomData(r)
	if err != nil {
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "join",
			Success:   false,
		})
		return
	}
	wsConn.Conn.WriteJSON(wsResponse{
		Operation: "join",
		Success:   true,
		Data:      d,
	})
	if len(r.Users) == 1 {
		wsConn.Session.SetAsLeader()
	}

}
func setRoomPlayback(rID string, wsConn *wsConnection, playing bool) {
	r, err := room.GetRoom(wsConn.Session.InRoom)
	if err == nil {
		r.SetPlaying(playing)
	}
}

func handleMessage(msg wsMessage, wsConn wsConnection) {
	switch msg.Operation {
	case "setUn":
		wsConn.Session.SetUsername(msg.Data.Username)
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "setUn",
			Success:   true,
		})
	case "create":
		room := room.CreateRoom(msg.Data.Room, wsConn.Session.ID)
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "create",
			Success:   true,
			Data:      room.RoomID,
		})

	case "sync":
		sync(msg.Data.Room, &wsConn, msg.Data.Time)
	case "change":
		changeRoomVideo(msg.Data.Room, &wsConn, msg.Data.VideoLink)
	case "join":
		joinRoom(msg.Data.Room, &wsConn)
	case "leave":
		wsConn.Session.LeaveRoom(msg.Data.Room)
	case "chat":
		// TODO
	case "leader":
		isLeader := wsConn.Session.IsLeader(msg.Data.Room)
		if isLeader {
			s, err := session.GetSess(msg.Data.UserID)
			if err != nil {
				s.SetAsLeader()
			}
			wsConn.Conn.WriteJSON(wsResponse{
				Operation: "leaderOff",
			})
		}

	case "hasUn":
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "hasUn",
			Data:      wsConn.Session.Username != "",
		})
	case "pause":
		setRoomPlayback(msg.Data.Room, &wsConn, false)
	case "play":
		setRoomPlayback(msg.Data.Room, &wsConn, true)
	}
}

func handleWS(c *gin.Context) {
	sess := session.AddSess()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println("Connection upgrade failed...")
	}
	wsConn := wsConnection{Conn: conn, Session: &sess}
	msg := wsMessage{}
	wsConn.Session.SetWsConn(conn)
	for {
		if wsConn.Conn.ReadJSON(&msg) != nil {
			wsConn.Session.Delete()
			break
		}
		handleMessage(msg, wsConn)
	}
}
