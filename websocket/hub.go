package websocket

import (
	"net/http"
	"sync"
	"time"

	"robotech/logger"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsClient struct {
	conn    *websocket.Conn
	sendCh  chan []byte
	session string
	closed  chan struct{}
}

var (
	hubMu       sync.Mutex
	subscribers = make(map[string]map[*wsClient]bool) // sessionID -> set of clients
)

func init() {
	// no-op for now
}

// HandleWS upgrades and registers client, starts read/write goroutines with heartbeat
func HandleWS(w http.ResponseWriter, r *http.Request, sessionID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("ws upgrade err:", err)
		return
	}
	client := &wsClient{
		conn:    conn,
		sendCh:  make(chan []byte, 256),
		session: sessionID,
		closed:  make(chan struct{}),
	}
	registerClient(sessionID, client)
	logger.Infof("ws client subscribed to session %s", sessionID)

	// reader - detect client close/pings
	go client.readPump()
	// writer + heartbeat
	go client.writePump()
}

func registerClient(sessionID string, c *wsClient) {
	hubMu.Lock()
	if subscribers[sessionID] == nil {
		subscribers[sessionID] = make(map[*wsClient]bool)
	}
	subscribers[sessionID][c] = true
	hubMu.Unlock()
}

// Push pushes a textual message to subscribers of sessionID
func Push(sessionID string, msg string) {
	hubMu.Lock()
	clients := subscribers[sessionID]
	var list []*wsClient
	for c := range clients {
		list = append(list, c)
	}
	hubMu.Unlock()
	// write outside lock
	for _, c := range list {
		select {
		case c.sendCh <- []byte(msg):
		default:
			// if send queue full, drop and remove client
			logger.Error("ws client send queue full, removing")
			c.close()
		}
	}
}

// CloseAll closes all clients for a session
func CloseAll(sessionID string) {
	hubMu.Lock()
	clients := subscribers[sessionID]
	delete(subscribers, sessionID)
	hubMu.Unlock()
	for c := range clients {
		c.close()
	}
}

// remove single client from subscribers
func removeClient(sessionID string, c *wsClient) {
	hubMu.Lock()
	if subscribers[sessionID] != nil {
		delete(subscribers[sessionID], c)
		if len(subscribers[sessionID]) == 0 {
			delete(subscribers, sessionID)
		}
	}
	hubMu.Unlock()
}

// wsClient methods

func (c *wsClient) readPump() {
	defer c.close()
	c.conn.SetReadLimit(1024 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			// reader ended
			return
		}
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.close()
	}()
	for {
		select {
		case <-c.closed:
			return
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Error("ws write err:", err)
				return
			}
		case <-ticker.C:
			// send ping
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				logger.Error("ws ping err:", err)
				return
			}
		}
	}
}

func (c *wsClient) close() {
	select {
	case <-c.closed:
		return
	default:
		close(c.closed)
		c.conn.Close()
		removeClient(c.session, c)
	}
}
