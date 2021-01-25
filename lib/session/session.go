package session

import (
	"context"

	"github.com/Aapeli123/watchwithfriends3/lib/database"
	"github.com/Aapeli123/watchwithfriends3/lib/room"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Username string
	ID       string
	InRoom   string
	wsConn   *websocket.Conn
}

func (s *Session) SetUsername(un string) {
	s.Username = un
	database.Sessions().Doc(s.ID).Set(context.Background(), s)
}

func (s *Session) Delete() {
	se, _ := GetSess(s.ID)
	if se.InRoom != "" {
		se.LeaveRoom(s.InRoom)
	}
	database.Sessions().Doc(se.ID).Delete(context.Background())
}

func AddSess() Session {
	sess := Session{ID: uuid.New().String()}
	database.Sessions().Doc(sess.ID).Set(context.Background(), sess)
	return sess
}

func GetSess(SID string) (Session, error) {
	snapshot, err := database.Sessions().Doc(SID).Get(context.Background())
	if err != nil {
		return Session{}, err
	}
	sess := Session{}
	snapshot.DataTo(&sess)
	return sess, nil
}

func ValidateSess(SID string) bool {
	_, err := GetSess(SID)
	if err != nil {
		return false
	}
	return true
}

func (s *Session) JoinRoom(rID string) {
	r, _ := room.GetRoom(rID)
	r.Users = append(r.Users, s.ID)
	r.UpdateRoom()
	s.InRoom = r.RoomID
	s.UpdateSession()

}

func (s *Session) SetWsConn(conn *websocket.Conn) {
	s.wsConn = conn
}

func (s *Session) SetAsLeader() {
	type msg struct {
		Operation string
	}
	if s.InRoom == "" {
		return
	}
	r, _ := room.GetRoom(s.InRoom)
	r.Leader = s.ID
	r.UpdateRoom()
	s.UpdateSession()
	s.wsConn.WriteJSON(msg{Operation: "leader"})
}

func (s *Session) LeaveRoom(rID string) {
	r, _ := room.GetRoom(rID)

	for i, id := range r.Users {
		if id == s.ID {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			break
		}
	}

	if len(r.Users) == 0 {
		r.Delete()
		return
	}
	if r.Leader == s.ID {
		r.Leader = r.Users[0]
	}
	s.InRoom = ""
	s.UpdateSession()
	r.UpdateRoom()

}

func (s *Session) IsLeader(rID string) bool {
	r, _ := room.GetRoom(rID)
	return r.Leader == s.ID
}

func (s *Session) UpdateSession() {
	database.Sessions().Doc(s.ID).Set(context.Background(), s)
}
