package websocket

import (
	"net/http"
	"robotech/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

var (
	subscribers   = make(map[string][]*websocket.Conn)
	subscribersMu sync.Mutex
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, sessionID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WebSocket upgrade failed:", err)
		return
	}

	subscribersMu.Lock()
	subscribers[sessionID] = append(subscribers[sessionID], conn)
	subscribersMu.Unlock()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				removeConn(sessionID, conn)
				return
			}
		}
	}()
	go heartbeat(conn)
}

func Push(sessionID, msg string) {
	subscribersMu.Lock()
	defer subscribersMu.Unlock()
	for _, conn := range subscribers[sessionID] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			conn.Close()
			removeConn(sessionID, conn)
		}
	}
}

func removeConn(sessionID string, target *websocket.Conn) {
	subscribersMu.Lock()
	defer subscribersMu.Unlock()
	list := subscribers[sessionID]
	newList := []*websocket.Conn{}
	for _, c := range list {
		if c != target {
			newList = append(newList, c)
		}
	}
	subscribers[sessionID] = newList
}

func CloseSubscribers(sessionID string) {
	subscribersMu.Lock()
	defer subscribersMu.Unlock()
	for _, c := range subscribers[sessionID] {
		c.Close()
	}
	delete(subscribers, sessionID)
}

func heartbeat(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			conn.Close()
			return
		}
	}
}
