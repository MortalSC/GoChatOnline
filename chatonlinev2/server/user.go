package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建用户的API
func NewUser(conn net.Conn, server *Server) *User {
	// 用户地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr, // 用户名默认用地址
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

// 用户上线业务处理
func (u *User) Online() {
	// 用户上线了！把用户加入到OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 对其他在线用户进行广播：上线通知
	u.server.BroadCast(u, "One of my friends is online")
}

// 用户下线业务处理
func (u *User) Offline() {
	// 用户下线了！把用户从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 对其他在线用户进行广播：上线通知
	u.server.BroadCast(u, "One of my friends is offline")
}

// 给当前User对应的客户端发消息
func (u *User) sendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户消息处理业务
func (u *User) DoMessage(msg string) {
	// u.server.BroadCast(u, msg)
	// 在线用户查询功能模拟
	if msg == "who" {

		// 打印当前在线用户数量
		// fmt.Println("当前在线用户数量:", len(u.server.OnlineMap))
		// 遇到who就查询当前在线的用户
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[ " + user.Addr + " ] " + user.Name + " : " + "online...\n"
			u.sendMsg(onlineMsg)

		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名的格式：rename|名称
		newName := msg[7:]
		// 判断name是否存在，不允许重名
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.sendMsg("The user name already exists!\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.sendMsg("Your username has been updated to: " + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊消息格式：to|username|msg

		// 1. 获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			// 用户名为空
			u.sendMsg("The message format is incorrect, please input like : [to|uesrname|msg]\n")
			return
		}

		// 2. 根据用户名得到对方的user对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			// 用户不存在
			u.sendMsg("The user is nil\n")
			return
		}

		// 3. 获取消息内容，并发送
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.sendMsg("The content you sent is empty\n")
			return
		}
		remoteUser.sendMsg("[" + u.Name + "] send the msg to you: [ " + content + " ]\n")

	} else {
		// 不是用户查询就作为广播消息
		u.server.BroadCast(u, msg)
	}
}
