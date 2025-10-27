package server

import (
	"fmt"
	"net"
	"robotech/logger"
	"robotech/session"
	"sync"
)

var (
	serverMu sync.Mutex
	servers  = make(map[string]net.Listener) // key: addr "ip:port"
)

func init() {
	// no-op for now
}

// StartTCPServer start listening on ip:port and accept connections.
// addrKey is ip:port
func StartTCPServer(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	serverMu.Lock()
	if _, ok := servers[addr]; ok {
		serverMu.Unlock()
		return nil // already running
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		serverMu.Unlock()
		return err
	}
	servers[addr] = ln
	serverMu.Unlock()

	logger.Infof("TCP server started on %s", addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logger.Error("accept error:", err)
				// if listener closed, exit
				serverMu.Lock()
				_, exists := servers[addr]
				serverMu.Unlock()
				if !exists {
					return
				}
				continue
			}
			go session.HandleIncoming(conn)
		}
	}()
	return nil
}

func StopTCPServer(addr string) {
	serverMu.Lock()
	ln, ok := servers[addr]
	if ok {
		_ = ln.Close()
		delete(servers, addr)
		logger.Infof("stopped server %s", addr)
	}
	serverMu.Unlock()
}

// Stop all (helper)
func StopAllServers() {
	serverMu.Lock()
	for addr, ln := range servers {
		_ = ln.Close()
		delete(servers, addr)
		logger.Infof("stopped server %s", addr)
	}
	serverMu.Unlock()
}
