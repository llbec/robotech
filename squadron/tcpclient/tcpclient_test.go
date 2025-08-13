package tcpclient_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	// Create a log channel for the client
	logChan := make(chan string, 1024*8)

	go func() {
		for logMsg := range logChan {
			t.Log(logMsg) // Log messages for testing
		}
	}()

	host := "192.168.1.104:8000"

	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		line, err := reader.ReadString('\n')
		if err != nil {
			logChan <- fmt.Sprintf("[Client] 读取失败: %v", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		logChan <- "[Client 收到] " + line
		break
	}
}

func TestClientRcv(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.1.104:8000")
	if err != nil {
		t.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()
	t.Log("已连接服务器，开始接收...")

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		line, err := reader.ReadString('\n')
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				// 超时后继续尝试读取，不退出循环
				continue
			}
			t.Fatal("读取错误:", err)
			break
		}
		t.Logf("接收: %s", line)
	}

	t.Log("连接结束")
}
