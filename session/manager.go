package session

import (
	"net"
	"robotech/logger"
	"robotech/websocket"
	"sync"
	"time"
)

type Session struct {
	ID       string
	Conn     net.Conn
	CreateAt time.Time
}

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.Mutex
)

func HandleConnection(conn net.Conn) {
	id := conn.RemoteAddr().String()
	s := &Session{ID: id, Conn: conn, CreateAt: time.Now()}
	sessionsMu.Lock()
	sessions[id] = s
	sessionsMu.Unlock()
	logger.Info("New session:", id)

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			logger.Error("Connection closed:", id)
			conn.Close()
			Close(id)
			return
		}
		msg := string(buf[:n])
		logger.Info("Recv:", msg)
		websocket.Push(id, msg)
	}
}

func RegisterClient(conn net.Conn) {
	id := conn.RemoteAddr().String()
	s := &Session{ID: id, Conn: conn, CreateAt: time.Now()}
	sessionsMu.Lock()
	sessions[id] = s
	sessionsMu.Unlock()
	logger.Info("Client connected:", id)
}

func List() []*Session {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	out := []*Session{}
	for _, s := range sessions {
		out = append(out, s)
	}
	return out
}

func Close(id string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	if s, ok := sessions[id]; ok {
		s.Conn.Close()
		delete(sessions, id)
		websocket.CloseSubscribers(id)
	}
}
