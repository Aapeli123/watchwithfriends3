package session

import (
	"context"

	"github.com/Aapeli123/watchwithfriends3/lib/database"
	"github.com/google/uuid"
)

type Session struct {
	Username string
	ID       string
	Active   bool
}

func (s *Session) SetUsername(un string) {
	s.Username = un
	database.Sessions().Doc(s.ID).Set(context.Background(), s)
}

func (s *Session) Delete() {
	database.Sessions().Doc(s.ID).Delete(context.Background())
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
