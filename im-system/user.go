package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建用户的API
func NewUser(conn net.Conn) *User {
	// 用户地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr, // 用户名默认用地址
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前User channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	// 不断监听
	for {
		msg := <-u.C
		// 如果有消息，就发送
		u.conn.Write([]byte(msg + "\n"))
	}
}
