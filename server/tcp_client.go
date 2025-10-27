package server

import (
	"fmt"
	"net"
	"robotech/logger"
	"robotech/session"
	"sync"
)

var (
	clientMu sync.Mutex
	clients  = make(map[string]net.Conn) // key: addr
)

// StartTCPClient dials to remote and registers the session
func StartTCPClient(ip string, port int) error {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	clientMu.Lock()
	if _, ok := clients[addr]; ok {
		clientMu.Unlock()
		return nil
	}
	clientMu.Unlock()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	clientMu.Lock()
	clients[addr] = conn
	clientMu.Unlock()

	session.RegisterOutgoing(conn)
	logger.Infof("TCP client connected to %s", addr)
	return nil
}

func StopTCPClient(addr string) {
	clientMu.Lock()
	if c, ok := clients[addr]; ok {
		_ = c.Close()
		delete(clients, addr)
		logger.Infof("stopped client %s", addr)
	}
	clientMu.Unlock()
}

func StopAllClients() {
	clientMu.Lock()
	for addr, c := range clients {
		_ = c.Close()
		delete(clients, addr)
		logger.Infof("stopped client %s", addr)
	}
	clientMu.Unlock()
}
