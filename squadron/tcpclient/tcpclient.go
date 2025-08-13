package tcpclient

import (
	"fmt"
	"io"
	"net"
	"robotech/armory"
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
	c.logChan <- fmt.Sprintf("[Client] 尝试连接服务器 %s ...", c.host)
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		c.logChan <- fmt.Sprintf("[Client] 连接失败: %v", err)
		return
	}
	c.conn = conn
	c.logChan <- fmt.Sprintf("[Client] 已连接服务器 %s", c.host)

	/*reader := bufio.NewReader(conn)

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
	}*/
	buf := make([]byte, 1024*8)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			//fmt.Print(string(buf[:n]))
			//fmt.Println(BytesToHex(buf[:n]))
			c.logChan <- "[Client 收到]\n" + armory.BytesToHex(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				//log.Println("Server closed the connection")
				c.logChan <- "[Client] Server closed the connection"
			} else {
				//log.Printf("Read error: %v\n", err)
				c.logChan <- fmt.Sprintf("[Client] 读取失败: %v", err)
			}
			break
		}
	}
	c.logChan <- "[Client] 连接关闭"
	conn.Close()
	c.conn = nil
	//c.logChan <- fmt.Sprintf("[Client] 连接到服务器 %s ", c.Host())
}

func (c *TCPClient) Send(data string) error {
	if c.conn == nil {
		return fmt.Errorf("未连接服务器")
	}
	_, err := c.conn.Write([]byte(data + "\n"))
	return err
}

func (c *TCPClient) SendHexString(hex string) error {
	if c.conn == nil {
		return fmt.Errorf("未连接服务器")
	}

	bytes, err := armory.HexToBytes(hex)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(append(bytes, '\n'))
	return err
}

func (c *TCPClient) IsConnected() bool {
	return c.conn != nil
}

func (c *TCPClient) Host() string {
	if c.conn == nil {
		return c.host + "(未连接)"
	}
	return c.conn.RemoteAddr().String()
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
