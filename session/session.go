package session

import (
	"errors"
	"github.com/google/uuid"
	"mq/academy/ent"
	"sync"
	"time"
)

type Session struct {
	ID        string    `json:"sessionID"`
	User      *ent.User `json:"user"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Manager interface {
	SessionKey(u *ent.User) string
	Set(u *ent.User) (*Session, error)
	Get(id string) (*Session, error)
}

type ManagerInMem struct {
	sync.Map
}

func NewManagerInMem() Manager {
	return &ManagerInMem{sync.Map{}}
}

func (mgr *ManagerInMem) SessionKey(u *ent.User) string {
	return uuid.New().String()
}

func (mgr *ManagerInMem) Set(u *ent.User) (*Session, error) {
	sid := mgr.SessionKey(u)
	sess := &Session{
		ID:        sid,
		User:      u,
		UpdatedAt: time.Now(),
	}
	mgr.Store(sid, sess)
	return sess, nil
}

func (mgr *ManagerInMem) Get(id string) (*Session, error) {
	x, ok := mgr.Load(id)
	if !ok {
		return nil, errors.New("invalid session key")
	}
	sess := x.(*Session)
	return sess, nil
}
