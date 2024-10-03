package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string           // ip
	Port      int              // 端口
	OnlineMap map[string]*User // 在线用户表
	mapLock   sync.RWMutex     // OnlineMap 是全局的，需要加锁
	Message   chan string      // 消息广播通道
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息的channel，一旦有消息就发送给全部的在线User
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		// 将msg发送给全部的在线User
		s.mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息API
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[ " + user.Addr + " ] " + user.Name + " : " + msg

	// 消息放入通道中
	s.Message <- sendMsg
}

// 业务处理
func (s *Server) Handler(conn net.Conn) {
	// fmt.Println("connect successful")

	// 新建用户
	user := NewUser(conn)
	// 用户上线了！把用户加入到OnlineMap中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// 对其他在线用户进行广播：上线通知
	s.BroadCast(user, "One of my friends is online\n")

	// 当前handler阻塞
	select {}
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

	// 启动对广播通道的内容监听
	go s.ListenMessager()

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
