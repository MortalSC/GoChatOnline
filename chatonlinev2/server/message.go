package main

// Message结构体，用于封装聊天消息
type Message struct {
	Sender  string `json:"sender"`  // 消息发送者
	Content string `json:"content"` // 消息内容
	Type    string `json:"type"`    // 消息类型，如"broadcast"（群聊）或"private"（私聊）
	Target  string `json:"target"`  // 目标用户名，用于私聊
}
