package wsservice

import (
	. "github.com/ataboo/gowssrv/models"
	"github.com/gorilla/websocket"
	"time"
	"fmt"
)

type connectedUser struct {
	user *User
	conn *websocket.Conn
	lastUpdate time.Time
}

type WsHub struct {
	users map[string]connectedUser
}

func (h *WsHub) Cleanup() {
	for _, user := range h.users {
		user.conn.Close()
	}
}

func (h *WsHub) AddUser(user *User, conn *websocket.Conn) error {
	id := user.ID.String()
	if h.HasUser(id) {
		h.RemoveUser(id)
	}

	h.users[user.ID.String()] = connectedUser{
		user,
		conn,
		time.Now(),
	}

	//TODO init?

	return nil
}

func (h *WsHub) RemoveUser(ID string) error {
	user, ok := h.users[ID]
	if !ok {
		return fmt.Errorf("failed to find user")
	}

	user.conn.Close()
	delete(h.users, ID)

	return nil
}

func (h *WsHub) HasUser(ID string) bool  {
	_, ok := h.users[ID]

	return ok
}

