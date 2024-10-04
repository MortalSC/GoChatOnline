# GoChatOnline
> 本项目是Go语言初学阶段性的项目实战《即时通信系统》





# 开发概述

## 1. `main.go`

- 这是程序的入口，启动服务端并监听特定IP和端口。
- 创建一个新的Server对象并调用`Start()`方法启动服务器。

```go
package main

func main() {
    server := NewServer("127.0.0.1", 8888) // 创建一个服务器实例，监听127.0.0.1的8888端口
    server.Start()                         // 启动服务器
}
```

## 2. `server.go`

- 定义了Server结构体，包含服务器的IP地址、端口、在线用户表、广播消息通道等。
- 服务器通过监听特定端口接收来自客户端的连接请求，并处理客户端的消息。
- 提供了消息广播、消息处理等功能。

主要方法：

- `NewServer`: 创建Server对象。
- `ListenMessager`: 监听广播消息并将其分发给所有在线用户。
- `BroadCast`: 广播消息API，将消息分发给所有在线用户。
- `Handler`: 处理客户端连接、消息接收和处理的逻辑。
- `Start`: 启动服务器，监听指定端口，处理客户端的连接。

```go
type Server struct {
    Ip        string
    Port      int
    OnlineMap map[string]*User
    mapLock   sync.RWMutex
    Message   chan string
}
```

## 3. `user.go`

- 定义了User结构体，代表一个在线用户，包括用户名、用户地址、连接信息等。
- 负责用户的上线、下线和消息处理。

主要方法：

- `NewUser`: 创建一个新的用户。
- `Online`: 处理用户上线逻辑。
- `Offline`: 处理用户下线逻辑。
- `DoMessage`: 处理用户发送的消息并根据命令格式进行相应操作（如私聊或广播）。

```go
type User struct {
    Name string
    Addr string
    C    chan string
    conn net.Conn
    server *Server
}
```

## 4. `client.go`

- 客户端实现文件，负责与服务器通信。
- 通过命令行界面，用户可以进行群聊、私聊、修改用户名等操作。

主要方法：

- `NewClient`: 创建客户端并连接服务器。
- `PublicChat`: 群聊模式，发送消息到所有用户。
- `PrivateChat`: 私聊模式，发送消息到指定用户。
- `UpdateName`: 更新用户名。
- `Run`: 负责处理客户端的运行逻辑，包括模式选择和业务逻辑处理。

```go
type Client struct {
    serverIp   string
    serverPort int
    Name       string
    conn       net.Conn
    flag       int
}
```





----

---



# 网络基础：相关调用

> Go语言内置了一个强大的`net`包，用于处理网络相关的操作。主要有以下几个关键的组件：

## **net.Listener**

`net.Listener`接口用于监听传入的网络连接，通常用来创建服务器。例如，可以用`net.Listen`函数来监听特定的IP和端口。

```go
listener, err := net.Listen("tcp", ":8080")
```

这个接口包含以下方法：

- `Accept()`: 接受传入的连接。
- `Close()`: 关闭监听器。
- `Addr()`: 返回监听器的网络地址。

## **net.Conn**

`net.Conn`接口表示一个网络连接，既可以是服务器端接收到的连接，也可以是客户端发起的连接。它是全双工的，支持同时读写数据。

- `Read()`: 从连接中读取数据。
- `Write()`: 向连接写入数据。
- `Close()`: 关闭连接。
- `LocalAddr()`: 获取本地地址。
- `RemoteAddr()`: 获取远程地址。

```go
conn, err := net.Dial("tcp", "example.com:80")
```

## **net.Dial 和 net.Listen**

这两个函数用于创建连接和监听连接：

- `net.Dial`：用于客户端发起网络连接，支持TCP、UDP、HTTP等协议。
- `net.Listen`：用于服务器端监听连接。

## **net.TCPConn / net.UDPConn**

这两个是具体的网络协议连接对象，分别表示TCP和UDP连接。

- `net.TCPConn`: 用于TCP连接，继承自`net.Conn`。
- `net.UDPConn`: 用于UDP连接，支持无连接的数据报传输。
