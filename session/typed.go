package session

import (
	"net"
	"sync"
	"time"
)

type MsgType string

const (
	MsgTypeString MsgType = "string"
	MsgTypeHex    MsgType = "hex"
)

type Session struct {
	ID        string
	Conn      net.Conn
	CreatedAt time.Time
	MsgType   MsgType

	SendChan chan []byte
	RecvChan chan string

	mu          sync.Mutex
	closed      bool
	assignedMap map[int]bool // actionID -> true
	// stop chan optionally for runners (runners managed in action package)
}

var (
	sessions    = make(map[string]*Session)
	sessionsMu  sync.Mutex
	maxRecvSize = 4096
)
