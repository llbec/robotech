package server

import (
	"fmt"
	"net"
	"robotech/logger"
	"robotech/session"
)

var servers = make(map[string]net.Listener)
var clients = make(map[string]net.Conn)

func StartTCPServer(ip string, port int) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("TCP Server start error:", err)
		return
	}
	servers[addr] = ln
	logger.Info("TCP Server listening:", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("Accept error:", err)
			return
		}
		go session.HandleConnection(conn)
	}
}

func StopAllServers() {
	for addr, ln := range servers {
		ln.Close()
		logger.Info("Closed server:", addr)
	}
}

func StartTCPClient(ip string, port int) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Error("TCP Client connect error:", err)
		return
	}
	clients[addr] = conn
	session.RegisterClient(conn)
}

func StopAllClients() {
	for addr, conn := range clients {
		conn.Close()
		logger.Info("Closed client:", addr)
	}
}
