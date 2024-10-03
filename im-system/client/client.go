package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// 创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	clinet := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 连接 Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}

	clinet.conn = conn

	// 返回对象
	return clinet
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>>>>>>> connect Server failed...")
		return
	}
	fmt.Println(">>>>>>>>>>>>>>> connect Server successful...")

	// 启动客户端业务
	//select {}
	for {

	}
}
