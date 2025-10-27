package session

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"robotech/logger"
	"robotech/utils"
	"robotech/websocket"
)

// Init (call at startup)
func init() {
	// no-op for now
}

// HandleIncoming called by server when an accepted connection arrives
func HandleIncoming(conn net.Conn) {
	id := conn.RemoteAddr().String()
	s := newSession(id, conn)
	addSession(s)
	logger.Infof("session new (accepted): %s", id)
}

// RegisterOutgoing used by client after successful Dial
func RegisterOutgoing(conn net.Conn) {
	id := conn.RemoteAddr().String()
	s := newSession(id, conn)
	addSession(s)
	logger.Infof("session new (outgoing): %s", id)
}

func newSession(id string, conn net.Conn) *Session {
	s := &Session{
		ID:          id,
		Conn:        conn,
		CreatedAt:   time.Now(),
		MsgType:     MsgTypeString,
		SendChan:    make(chan []byte, 128),
		RecvChan:    make(chan string, maxRecvSize),
		assignedMap: make(map[int]bool),
	}
	// start io loops
	go s.readLoop()
	go s.writeLoop()
	return s
}

func addSession(s *Session) {
	sessionsMu.Lock()
	sessions[s.ID] = s
	sessionsMu.Unlock()
	// notify websockets about open
	websocket.Push(s.ID, "[EVENT] session_open")
}

// List sessions snapshot
func List() []map[string]any {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	out := make([]map[string]any, 0, len(sessions))
	for id, s := range sessions {
		s.mu.Lock()
		acts := make([]int, 0, len(s.assignedMap))
		for k := range s.assignedMap {
			acts = append(acts, k)
		}
		s.mu.Unlock()
		out = append(out, map[string]any{
			"id":         id,
			"remote":     s.Conn.RemoteAddr().String(),
			"local":      s.Conn.LocalAddr().String(),
			"created_at": s.CreatedAt.Format(time.RFC3339),
			"actions":    acts,
		})
	}
	return out
}

func Get(id string) *Session {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	return sessions[id]
}

func Close(id string) {
	sessionsMu.Lock()
	s, ok := sessions[id]
	if ok {
		delete(sessions, id)
	}
	sessionsMu.Unlock()
	if ok {
		s.close()
	}
}

// internal close
func (s *Session) close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	s.mu.Unlock()

	_ = s.Conn.Close()
	close(s.SendChan)
	close(s.RecvChan)
	websocket.CloseAll(s.ID)
	logger.Infof("session closed: %s", s.ID)
}

// readLoop reads raw bytes and publishes to RecvChan and websocket
func (s *Session) readLoop() {
	reader := bufio.NewReader(s.Conn)
	buf := make([]byte, maxRecvSize)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Error("read error for session", s.ID, err)
			}
			// close session and exit
			s.close()
			return
		}
		if n <= 0 {
			continue
		}
		data := buf[:n]

		txt := string(data)
		if s.MsgType == MsgTypeHex {
			txt = utils.BytesToHexString(data)
		}
		logger.Infof("%s recv: %s", s.ID, txt)
		// push hex and text to websocket
		websocket.Push(s.ID, "[RECV_TEXT] "+txt)

		// send to RecvChan (string form for action matching)
		select {
		case s.RecvChan <- txt:
		default:
			// drop if full
			logger.Error("recv chan full for session", s.ID)
		}
	}
}

// writeLoop writes []byte to conn; messages on SendChan are raw bytes
func (s *Session) writeLoop() {
	for b := range s.SendChan {
		_, err := s.Conn.Write(b)
		if err != nil {
			logger.Error("write error for session", s.ID, err)
			s.close()
			return
		}

		txt := string(b)
		if s.MsgType == MsgTypeHex {
			txt = utils.BytesToHexString(b)
		}
		// push to ws as hex and text
		websocket.Push(s.ID, "[SENT_TEXT] "+txt)
		logger.Infof("%s sent: %s", s.ID, txt)
	}
}

// GetSessionMsgType returns message type for session
func GetSessionMsgType(sessionID string) MsgType {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s := sessions[sessionID]
	if s == nil {
		return MsgTypeString
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.MsgType
}

// SetSessionMsgType sets message type for session
// MsgTypeHex: convert to hex string before sending
// MsgTypeString: send as is
func SetSessionMsgType(sessionID string, msgType MsgType) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s := sessions[sessionID]
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MsgType = msgType
	logger.Infof("session %s set msg type to %s", s.ID, msgType)
}

// Send sends raw bytes to session
func Send(sessionID string, b []byte) error {
	s := Get(sessionID)
	if s == nil {
		return fmt.Errorf("session not found")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("session closed")
	}
	select {
	case s.SendChan <- b:
		return nil
	default:
		return fmt.Errorf("send queue full")
	}
}

// AssignAction mark action bound to session (session side)
func AssignAction(sessionID string, actionID int) error {
	s := Get(sessionID)
	if s == nil {
		return fmt.Errorf("session not found")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.assignedMap == nil {
		s.assignedMap = make(map[int]bool)
	}
	if s.assignedMap[actionID] {
		return fmt.Errorf("already assigned")
	}
	s.assignedMap[actionID] = true
	return nil
}

// UnassignAction remove assignment
func UnassignAction(sessionID string, actionID int) {
	s := Get(sessionID)
	if s == nil {
		return
	}
	s.mu.Lock()
	delete(s.assignedMap, actionID)
	s.mu.Unlock()
}

// IsActionAssigned returns whether action bound
func IsActionAssigned(sessionID string, actionID int) bool {
	s := Get(sessionID)
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.assignedMap[actionID]
}
