package models

import (
	"log"
	"math/rand"
	"sync"
)


type RoomMap struct {
  Mutex sync.RWMutex
  Map   map[string][]*Participant
}


func (r *RoomMap) Init() {
	r.Map = make(map[string][]*Participant)
}


func (r *RoomMap) Get(roomID string) []*Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomID := string(b)

	r.Map[roomID] = []*Participant{}

	return roomID
}

func (r *RoomMap) InsertIntoRoom(roomID string, participant *Participant) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], participant)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)

}
