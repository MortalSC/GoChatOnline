package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	mode       int // 当前客户端选择的模式：1.群聊；2.私聊；3.修改用户名
}

// 创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 连接 Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}

	// 创建客户端对象
	clinet := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		conn:       conn,
		mode:       0,
	}

	// 返回对象
	return clinet
}

// 展示客户端功能
func (client *Client) menu() bool {
	var mode int
	fmt.Println("1. Group chat mode")
	fmt.Println("2. Private chat mode")
	fmt.Println("3. Update user name")
	fmt.Println("0. Quit")

	fmt.Scanln(&mode)
	if mode >= 0 && mode <= 3 {
		client.mode = mode
		return true
	} else {
		fmt.Println("Invalid option, please choose again.")
		return false
	}
}

// 发送消息到服务器
func (client *Client) SendMessage(content string) {
	msg := Message{
		Sender:  client.Name, // 消息发送者
		Content: content,     // 消息内容
		Type:    "broadcast", // 默认消息为群聊
	}

	// 消息编码为json格式
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err : ", err)
		return
	}

	client.conn.Write(data) // 通过连接发送消息
}

// 群聊业务
func (client *Client) GroupChat() {

	fmt.Println("Enter your message (type 'exit' to quit):")
	reader := bufio.NewReader(os.Stdin)

	for {
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)

		if content == "exit" {
			break
		}

		// 创建群聊消息
		msg := Message{
			Sender:  client.Name,
			Content: content,
			Type:    "broadcast",
		}

		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("json.Marshal err : ", err)
			return
		}

		// 发送消息到服务器
		client.conn.Write(data)
	}
}

// 私聊业务
func (client *Client) PrivateChat() {
	fmt.Println("Enter the username you want to chat with:")
	reader := bufio.NewReader(os.Stdin)
	targetUser, _ := reader.ReadString('\n')
	targetUser = strings.TrimSpace(targetUser)

	fmt.Printf("You are now chatting privately with %s (type 'exit' to quit):\n", targetUser)

	for {
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)

		if content == "exit" {
			break
		}

		// 创建一个私聊消息
		msg := Message{
			Sender:  client.Name,
			Content: content,
			Type:    "private",  // 消息类型为私聊
			Target:  targetUser, // 私聊的目标用户
		}

		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("json.Marshal err : ", err)
			return
		}

		client.conn.Write(data) // 发送消息到服务器
	}
}

// 更新用户名
func (client *Client) UpdateUsername() {
	fmt.Println("Enter your new username:")
	reader := bufio.NewReader(os.Stdin)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	if newName == "" {
		fmt.Println("Username cannot be empty.")
		return
	}

	client.Name = newName

	msg := Message{
		Sender:  "System",
		Content: client.Name,
		Type:    "rename",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}

	_, err = client.conn.Write(data) // 发送更新用户名消息
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}

// 处理Server端回应的操作消息
// 比如修改用户成功后，提示的：修改成功等内容
func (client *Client) DealResponse() {
	// 一旦client有数据，就拷贝到标准输出，永久阻塞
	//io.Copy(os.Stdout, client.conn)

	reader := bufio.NewReader(client.conn)
	for {
		// 读取从服务器接收到的消息，按行读取
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error : ", err)
			return
		}

		// 去除读取到的换行符
		msg = strings.TrimSpace(msg)

		// 解析 json
		var message Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			fmt.Println("Unmarshal error : ", err)
			continue
		}

		// 根据消息的类型和内容显示
		switch message.Type {
		case "broadcast":
			fmt.Printf("[Broadcast] %s: %s\n", message.Sender, message.Content)
		case "private":
			fmt.Printf("[Private] %s: %s\n", message.Sender, message.Content)
		case "info":
			fmt.Printf("[Info] %s\n", message.Content)
		default:
			fmt.Printf("[Unknown] %s: %s\n", message.Sender, message.Content)
		}
	}
}

// 业务处理
func (client *Client) Run() {
	for {
		if client.menu() {

			// 根据不同的模式，进行不同的业务
			switch client.mode {
			case 1: // 群聊
				//fmt.Println("select : Group chat mode...")
				client.GroupChat()
			case 2: // 私聊
				//fmt.Println("select : Private chat mode...")
				client.PrivateChat()
			case 3: // 修改用户名
				//fmt.Println("select : Update user name...")
				client.UpdateUsername()
			case 0:
				fmt.Println("Exiting...")
				return // 退出
			}
		}
	}
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
		fmt.Println("Connection failed...")
		return
	}

	// 处理Server回应的操作提示消息
	go client.DealResponse()

	fmt.Println("Connection Server successful...")

	// 启动客户端业务
	client.Run()
}
