package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type TCPServer struct {
	addr    string
	ln      net.Listener
	logChan chan string

	conn net.Conn
}

func NewTCPServer(addr string, logChan chan string) *TCPServer {
	return &TCPServer{
		addr:    addr,
		logChan: logChan,
	}
}

func (s *TCPServer) Start() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logChan <- fmt.Sprintf("[Server] 监听失败: %v", err)
		return
	}
	s.ln = ln
	s.logChan <- fmt.Sprintf("[Server] 监听地址: %s", s.addr)

	conn, err := ln.Accept()
	if err != nil {
		s.logChan <- fmt.Sprintf("[Server] Accept失败: %v", err)
		return
	}
	s.conn = conn
	s.logChan <- "[Server] 客户端已连接"

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		line, err := reader.ReadString('\n')
		if err != nil {
			s.logChan <- fmt.Sprintf("[Server] 读取失败: %v", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s.logChan <- "[Server 收到] " + line
	}
	s.logChan <- "[Server] 连接关闭"
}

func (s *TCPServer) Send(data string) error {
	if s.conn == nil {
		return fmt.Errorf("没有客户端连接")
	}
	_, err := s.conn.Write([]byte(data + "\n"))
	return err
}

func (s *TCPServer) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.ln != nil {
		s.ln.Close()
	}
}
