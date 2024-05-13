package data

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

var SetChannel = make(chan SetMessage)
var DeleteChannel = make(chan DeleteMessage)

type SetMessage struct {
	ID   uuid.UUID
	Body []byte
}

type DeleteMessage struct {
	ID uuid.UUID
}

type Data struct {
	Data  map[uuid.UUID][]byte `json:"data"`
	Mutex sync.Mutex
}

func (db *Data) GetById(id uuid.UUID) (val []byte, err error) {
	db.Mutex.Lock()
	val, ok := db.Data[id]
	db.Mutex.Unlock()
	if !ok {
		return nil, errors.New("not found")
	}
	return val, nil
}

func (db *Data) Set(id uuid.UUID, bytes []byte) {
	SetChannel <- SetMessage{
		ID:   id,
		Body: bytes,
	}
}

func (db *Data) Delete(id uuid.UUID) {
	DeleteChannel <- DeleteMessage{
		ID: id,
	}
}
