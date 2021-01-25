package room

import (
	"context"
	"math/rand"

	"github.com/Aapeli123/watchwithfriends3/lib/database"
	"github.com/google/uuid"
)

// Needed for room code generation
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const roomCodeLen = 4

func createRoomcode() string {
	code := []byte{}
	for i := 0; i < roomCodeLen; i++ {
		r := rand.Intn(len(charset) - 1)
		char := charset[r]
		code = append(code, char)
	}
	return string(code)
}

type Room struct {
	RoomCode     string
	RoomID       string
	VideoPos     float64
	VideoID      string
	RoomName     string
	Users        []string
	Leader       string
	ChatMessages []string
	Playing      bool
}

func (r *Room) UpdateVideoTime(time float64) {
	r.VideoPos = time
	r.UpdateRoom()
}

func (r *Room) UpdateVideoLink(link string) {
	r.VideoID = link
	r.UpdateRoom()
}

func CreateRoom(roomName string, creator string) Room {
	room := Room{
		RoomName: roomName,
		VideoPos: 0,
		VideoID:  "",
		Leader:   creator,
		Users:    []string{},
		RoomCode: createRoomcode(),
		RoomID:   uuid.New().String(),
		Playing:  false,
	}
	database.Rooms().Doc(room.RoomID).Set(context.Background(), room)
	return room
}

func GetRoom(id string) (Room, error) {
	snapshot, err := database.Rooms().Doc(id).Get(context.Background())
	if err != nil {
		return Room{}, err
	}
	r := Room{}
	err = snapshot.DataTo(&r)
	if err != nil {
		return Room{}, err
	}
	return r, nil
}
func GetRooms() ([]Room, error) {
	snapshots, err := database.Rooms().Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	var rooms []Room
	var room Room
	for _, snapshot := range snapshots {
		snapshot.DataTo(&room)
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *Room) SetPlaying(playing bool) {
	r.Playing = playing
	r.UpdateRoom()
}

func (r *Room) Delete() {
	database.Rooms().Doc(r.RoomID).Delete(context.Background())
}

func (r *Room) UpdateRoom() {
	database.Rooms().Doc(r.RoomID).Set(context.Background(), r)
}
