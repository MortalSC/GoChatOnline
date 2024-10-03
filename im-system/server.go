package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string // ip
	Port int    // 端口
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// 业务处理
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("connect successful")
}

// 启动服务器的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err : ", err)
		return
	}
	// 利用语法特点，在这直接设置延迟关闭
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept() // 会阻塞
		if err != nil {
			fmt.Println("listener accept err : ", err)
			continue
		}

		// do handle
		go s.Handler(conn)
	}

	// close socket
}
