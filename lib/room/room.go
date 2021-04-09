package room

import (
	"errors"
	"math/rand"
)

// Needed for room code generation
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const roomCodeLen = 4

var ErrRoomNotFound = errors.New("Room not found")

var Rooms = map[string]*Room{}

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
}

func (r *Room) UpdateVideoLink(link string) {
	r.VideoID = link
}

func CreateRoom(roomName string, creator string) *Room {
	room := Room{
		RoomName: roomName,
		VideoPos: 0,
		VideoID:  "",
		Leader:   creator,
		Users:    []string{},
		RoomCode: createRoomcode(),
		Playing:  false,
	}
	Rooms[room.RoomCode] = &room
	return &room
}

func GetRoom(code string) (*Room, error) {
	if Rooms[code] == nil {
		return nil, ErrRoomNotFound
	}
	return Rooms[code], nil
}
func GetRooms() map[string]*Room {
	return Rooms
}
func RoomArray() []*Room {
	rooms := []*Room{}
	for _, r := range Rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

func (r *Room) SetPlaying(playing bool) {
	r.Playing = playing
}

func (r *Room) Delete() {
	delete(Rooms, r.RoomCode)
}
