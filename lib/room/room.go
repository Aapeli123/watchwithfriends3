package room

import session "github.com/Aapeli123/watchwithfriends3/lib/Session"

type Room struct {
	RoomCode string
	RoomID   string
	VideoPos int64
	VideoID  string
	RoomName string
	Users    []*session.Session
	Leader   *session.Session
}

func (r *Room) AddUser(u *session.Session) {
	r.Users = append(r.Users, u)
}
