package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, s)

	user.Online()

	// 用户活跃监听
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// 如果读取到的消息长度为 0， 表示用户关闭
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Connect Read err : ", err)
				return
			}

			// 提取用户的消息，去除 '\n'
			msg := string(buf[:n-1])

			// 将得到的消息广播
			// s.BroadCast(user, msg)
			user.DoMessage(msg)

			// 用户触发消息就认为是活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，重置计时器
		case <-time.After(time.Second * 300): // 5分钟
			// 超时
			// 将当前的User强制关闭
			user.sendMsg("You are offline\n")

			// 销毁占用的资源
			close(user.C)

			// 关闭连接
			conn.Close()

			// 退出当前的handler
			return
		}
	}
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
