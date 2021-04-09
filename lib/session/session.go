package session

import (
	"errors"

	"github.com/Aapeli123/watchwithfriends3/lib/room"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Sessions = map[string]*Session{}
var ErrNotInRoom = errors.New("User not in room")
var ErrAlreadyInRoom = errors.New("User is already in this room")

type Session struct {
	Username string
	ID       string
	InRoom   string
	wsConn   *websocket.Conn
}

func (s *Session) SetUsername(un string) {
	s.Username = un
}

func (s *Session) Delete() {
	if s.InRoom != "" {
		s.LeaveRoom(s.InRoom)
	}
	delete(Sessions, s.ID)
}

func AddSess() *Session {
	sess := Session{ID: uuid.New().String()}
	Sessions[sess.ID] = &sess
	return &sess
}

func GetSess(SID string) *Session {
	sess := Sessions[SID]
	return sess
}

func ValidateSess(SID string) bool {
	return true
}

func (s *Session) JoinRoom(rID string) error {
	r, err := room.GetRoom(rID)
	if err != nil {
		return err
	}
	if r.RoomCode == s.InRoom {
		return ErrAlreadyInRoom
	}

	r.Users = append(r.Users, s.ID)
	s.InRoom = r.RoomCode
	return nil
}

func (s *Session) SetWsConn(conn *websocket.Conn) {
	s.wsConn = conn
}

func (s *Session) SetAsLeader() error {
	type msg struct {
		Operation string
	}
	if s.InRoom == "" {
		return ErrNotInRoom
	}
	r, err := room.GetRoom(s.InRoom)
	if err != nil {
		return err
	}
	r.Leader = s.ID
	s.wsConn.WriteJSON(msg{Operation: "leader"})
	return nil
}

func (s *Session) RecommendVideo(videoID string) error {

	r, err := room.GetRoom(s.InRoom)
	if err != nil {
		return err
	}
	roomLeader := GetSess(r.Leader)
	type VideoRecommendation struct {
		Recommender string
		VideoID     string
	}
	type msg struct {
		Operation string
		Data      interface{}
	}
	roomLeader.wsConn.WriteJSON(msg{
		Operation: "recommendation",
		Data: VideoRecommendation{
			Recommender: s.Username,
			VideoID:     videoID,
		},
	})
	return nil
}

func (s *Session) LeaveRoom(rID string) error {

	r, err := room.GetRoom(rID)
	if err != nil {
		return ErrNotInRoom
	}
	for i, id := range r.Users {
		if id == s.ID {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			break
		}
	}

	if len(r.Users) == 0 {
		s.InRoom = ""
		r.Delete()
		return nil
	}
	if r.Leader == s.ID {
		r.Leader = r.Users[0]
		GetSess(r.Leader).SetAsLeader()
	}
	s.InRoom = ""
	type userData struct {
		Username string
		UserID   string
	}
	type roomData struct {
		room.Room
		Usernames []userData
	}
	type msg struct {
		Operation string
		Data      []userData
	}
	var rd roomData
	for _, user := range r.Users {
		s := GetSess(user)
		rd.Usernames = append(rd.Usernames, userData{Username: s.Username, UserID: s.ID})
	}
	for _, u := range r.Users {
		sess := GetSess(u)
		sess.SendToClient(msg{Operation: "userDisconnected", Data: rd.Usernames})
	}
	return nil
}

func (s *Session) IsLeader(rCode string) (bool, error) {
	r, err := room.GetRoom(rCode)
	if err != nil {
		return false, ErrNotInRoom
	}
	return r.Leader == s.ID, nil
}

func (s *Session) SendPlaybackState(play bool) {
	type msg struct {
		Operation string
	}
	if play {
		s.wsConn.WriteJSON(msg{
			Operation: "play",
		})
		return
	}
	s.wsConn.WriteJSON(msg{
		Operation: "pause",
	})
}

func (s *Session) SendToClient(msg interface{}) {
	if s.wsConn == nil {
		return
	}
	s.wsConn.WriteJSON(msg)
}
