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

func sync(wsConn *wsConnection, time float64) {
	r, err := room.GetRoom(wsConn.Session.InRoom)
	if err != nil {
		return
	}
	isLeader, err := wsConn.Session.IsLeader(r.RoomCode)
	if err != nil {
		return
	}
	if isLeader {
		r.UpdateVideoTime(time)
	}
	wsConn.Conn.WriteJSON(wsResponse{
		Operation: "sync",
		Data:      r.VideoPos,
		Success:   true,
	})

}

func changeRoomVideo(wsConn *wsConnection, vidLink string) {
	r, err := room.GetRoom(wsConn.Session.InRoom)
	if err != nil {
		return
	}
	r.UpdateVideoLink(vidLink)
	sendToAllInRoom(r.RoomCode, wsResponse{
		Operation: "change",
		Data:      vidLink,
	})
}

type roomData struct {
	room.Room
	Usernames []userData
}
type userData struct {
	Username string
	UserID   string
}

func getRoomData(r *room.Room) (roomData, error) {
	data := roomData{
		Room:      *r,
		Usernames: []userData{},
	}
	for _, u := range r.Users {
		s := session.GetSess(u)
		un := s.Username
		data.Usernames = append(data.Usernames, userData{
			Username: un,
			UserID:   s.ID,
		})
	}
	return data, nil
}

func sendToAllInRoom(rID string, msg interface{}) {
	room, err := room.GetRoom(rID)
	if err != nil {
		return
	}
	for _, u := range room.Users {
		sess := session.GetSess(u)
		sess.SendToClient(msg)
	}
}

func joinRoom(rID string, wsConn *wsConnection) {
	err := wsConn.Session.JoinRoom(rID)
	if err != nil {
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "join",
			Success:   false,
		})
		return
	}
	r, err := room.GetRoom(rID)
	if err != nil {
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "join",
			Success:   false,
		})
	}
	d, err := getRoomData(r)
	if err != nil {
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "join",
			Success:   false,
		})
		return
	}
	sendToAllInRoom(rID, wsResponse{
		Operation: "userConnected",
		Data:      d.Usernames,
	})
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
	if err != nil {
		return
	}
	if r.Playing == playing {
		return
	}
	r.SetPlaying(playing)
	for _, user := range r.Users {
		s := session.GetSess(user)
		s.SendPlaybackState(playing)
	}
}

func handleMessage(msg wsMessage, wsConn wsConnection) {
	switch msg.Operation {
	case "setUn":
		wsConn.Session.SetUsername(msg.Data.Username)
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "setUn",
			Success:   true,
			Data:      msg.Data.Username,
		})
	case "create":
		room := room.CreateRoom(msg.Data.Room, wsConn.Session.ID)
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "create",
			Success:   true,
			Data:      room.RoomCode,
		})
	case "rooms":
		rooms := room.RoomArray()
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "rooms",
			Success:   true,
			Data:      rooms,
		})
		break
	case "sync":
		sync(&wsConn, msg.Data.Time)
	case "change":
		changeRoomVideo(&wsConn, msg.Data.VideoLink)
	case "join":
		joinRoom(msg.Data.Room, &wsConn)
	case "leave":
		wsConn.Session.LeaveRoom(msg.Data.Room)
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "leave",
			Success:   true,
		})
	case "uid":
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "uid",
			Success:   true,
			Data:      wsConn.Session.ID,
		})
	case "chat":
		type chatMsg struct {
			Sender  string
			Content string
		}
		s := wsConn.Session
		sendToAllInRoom(s.InRoom, wsResponse{
			Operation: "chat",
			Data:      chatMsg{Sender: s.Username, Content: msg.Data.ChatMessage},
		})
	case "leader":
		isLeader, err := wsConn.Session.IsLeader(msg.Data.Room)
		if err != nil {
			break
		}
		if isLeader {
			s := session.GetSess(msg.Data.UserID)
			s.SetAsLeader()
			wsConn.Conn.WriteJSON(wsResponse{
				Operation: "leaderOff",
			})
		}

	case "hasUn":
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "hasUn",
			Success:   wsConn.Session.Username != "",
			Data:      wsConn.Session.Username,
		})
	case "recommendation":
		wsConn.Session.RecommendVideo(msg.Data.VideoLink)
	case "pause":
		setRoomPlayback(msg.Data.Room, &wsConn, false)
	case "play":
		setRoomPlayback(msg.Data.Room, &wsConn, true)
	case "ping":
		wsConn.Conn.WriteJSON(wsResponse{
			Operation: "pong",
		})
	}
}

func handleWS(c *gin.Context) {
	sess := session.AddSess()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println("Connection upgrade failed...")
	}
	wsConn := wsConnection{Conn: conn, Session: sess}
	msg := wsMessage{}
	wsConn.Session.SetWsConn(conn)
	for {
		if wsConn.Conn.ReadJSON(&msg) != nil {
			wsConn.Session.LeaveRoom(wsConn.Session.InRoom)
			wsConn.Session.Delete()
			break
		}
		handleMessage(msg, wsConn)
	}
}
