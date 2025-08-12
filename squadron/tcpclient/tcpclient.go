package tcpclient

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type TCPClient struct {
	host    string // TCP server host ip:port
	conn    net.Conn
	logChan chan string
}

func NewTCPClient(host string, logChan chan string) *TCPClient {
	return &TCPClient{
		host:    host,
		logChan: logChan,
		conn:    nil,
	}
}

func (c *TCPClient) Start() {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		c.logChan <- fmt.Sprintf("[Client] 连接失败: %v", err)
		return
	}
	defer conn.Close()
	c.conn = conn
	c.logChan <- fmt.Sprintf("[Client] 已连接服务器 %s", c.host)

	reader := bufio.NewReader(conn)

	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		line, err := reader.ReadString('\n')
		if err != nil {
			c.logChan <- fmt.Sprintf("[Client] 读取失败: %v", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		c.logChan <- "[Client 收到] " + line
	}
	c.logChan <- "[Client] 连接关闭"
}

func (c *TCPClient) Send(data string) error {
	if c.conn == nil {
		return fmt.Errorf("未连接服务器")
	}
	_, err := c.conn.Write([]byte(data + "\n"))
	return err
}

func (c *TCPClient) IsConnected() bool {
	return c.conn != nil
}

func (c *TCPClient) Host() string {
	if c.conn == nil {
		return c.host + "(未连接)"
	}
	return c.conn.RemoteAddr().Network()
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
