package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端选择的模式
}

// 创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	clinet := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
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

// 展示客户端功能
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. Group chat mode")
	fmt.Println("2. Private chat mode")
	fmt.Println("3. Update user name")
	fmt.Println("0. Quit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>> Invalid input <<<<<<")
		return false
	}
}

// 群聊业务
func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>> please input msg(to public)[exit means overchat] : ")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发送给服务器
		// 消息非空就发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err : ", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>> please input msg(to public)[exit means overchat] : ")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return
	}
}

// 私聊业务
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUser()
	fmt.Println(">>>>> please input name to private chat, [exit meas overchat] : ")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>> please input msg(to private)[exit means overchat] : ")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 发送给服务器
			// 消息非空就发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err : ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>> please input msg(to private)[exit means overchat] : ")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUser()
		fmt.Println(">>>>> please input name to private chat, [exit meas overchat] : ")
		fmt.Scanln(&remoteName)
	}
}

// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> please input new name :")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return false
	}

	return true
}

// 处理Server端回应的操作消息
// 比如修改用户成功后，提示的：修改成功等内容
func (client *Client) DealResponse() {
	// 一旦client有数据，就拷贝到标准输出，永久阻塞
	io.Copy(os.Stdout, client.conn)
}

// 业务处理
func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {

		}

		// 根据不同的模式，进行不同的业务
		switch client.flag {
		case 1: // 群聊
			//fmt.Println("select : Group chat mode...")
			client.PublicChat()
		case 2: // 私聊
			//fmt.Println("select : Private chat mode...")
			client.PrivateChat()
		case 3: // 修改用户名
			//fmt.Println("select : Update user name...")
			client.UpdateName()
		}
	}
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>>>>> connect Server failed...")
		return
	}

	// 处理Server回应的操作提示消息
	go client.DealResponse()

	fmt.Println(">>>>>>>>>>>>>>> connect Server successful...")

	// 启动客户端业务
	client.Run()
}
