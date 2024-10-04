package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// 定义Client接口，用于不同类型的客户端
type Client interface {
	SendMessage(msg Message) // 发送消息的接口
	GetName() string         // 获取客户端名称的接口
	Online()                 // 用户上线的接口
	Offline()                // 用户下线的接口
}

// User结构体实现了Client接口，代表一个用户
type User struct {
	Name   string       // 用户名
	Addr   string       // 用户地址
	C      chan Message // 用户的消息通道
	conn   net.Conn     // 用户的网络连接
	server *Server      // 用户所属的服务器
}

// 创建用户的API
func NewUser(conn net.Conn, server *Server) *User {
	// 用户地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr, // 用户名默认用地址
		Addr:   userAddr,
		C:      make(chan Message),
		conn:   conn,
		server: server,
	}

	// 启动监听当前用户消息通道的goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务处理
func (u *User) Online() {
	// 用户上线了！把用户加入到OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 对其他在线用户进行广播：上线通知
	u.server.BroadCast(u, "has joined")
}

// 用户下线业务处理
func (u *User) Offline() {
	// 用户下线了！把用户从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 对其他在线用户进行广播：上线通知
	u.server.BroadCast(u, "has left")
}

// 发消息给当前用户
func (u *User) SendMessage(msg Message) {
	// 将Message结构体转换成json格式
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal error : ", err)
		return
	}
	// 将json数据写入到用户连接中，发送给客户端
	//u.conn.Write(data)
	u.conn.Write(append(data, '\n'))
}

// 用户消息处理，根据消息类型进行处理
func (u *User) DoMessage(content string) {
	// 定义一个Message结构体
	var msg Message
	// json -> Message
	if err := json.Unmarshal([]byte(content), &msg); err != nil {
		fmt.Println("Message formamt error : ", err)
		//fmt.Println("Received content:", content) // 打印接收到的内容
		return
	}

	// 根据消息的Type执行不同的操作
	switch msg.Type {
	case "private": // 私聊消息
		// 查找用户
		u.server.mapLock.Lock()
		remoteUser, ok := u.server.OnlineMap[msg.Target]
		u.server.mapLock.Unlock()

		if !ok {
			// 如果目标用户不存在，发送错误消息
			u.SendMessage(
				Message{
					Sender:  "System",
					Content: fmt.Sprintf("User %s not found.\n", msg.Target),
					Type:    "error",
				})
			return
		}

		// 构建私聊消息
		privateMsg := Message{
			Sender:  u.Name,
			Content: msg.Content,
			Type:    "private",
		}

		// 消息发送给目标用户给
		remoteUser.SendMessage(privateMsg)

	case "rename": // 修改用户名
		u.server.mapLock.Lock()
		delete(u.server.OnlineMap, u.Name)
		u.Name = msg.Content
		u.server.OnlineMap[u.Name] = u
		u.server.mapLock.Unlock()

		u.SendMessage(
			Message{
				Sender:  "System",
				Content: "Username updated successfully.\n",
				Type:    "info",
			})

	default: // 默认群聊
		u.server.BroadCast(u, msg.Content)
	}
}

// 监听当前User消息管道的方法，一旦有消息，就发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C // 从消息管道中读取消息
		// Message -> json
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("json.Marshal error : ", err)
			return
		}
		// 将json数据写入到用户连接中，发送给客户端
		u.conn.Write(data)
	}
}

// 获取用户名称
func (u *User) GetName() string {
	return u.Name
}
