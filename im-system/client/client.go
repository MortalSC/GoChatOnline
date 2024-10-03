package main

import (
	"flag"
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

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip (defualt 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server port (defualt 8888)")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>>>>> connect Server failed...")
		return
	}
	fmt.Println(">>>>>>>>>>>>>>> connect Server successful...")

	// 启动客户端业务
	//select {}
}
