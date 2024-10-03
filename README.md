# GoChatOnline
> 本项目是Go语言初学阶段性的项目实战《即时通信系统》





# Server服务初步

## 接口展示

```go
// 创建一个server接口
func NewServer(ip string, port int) *Server { ... }

// 业务处理
func (s *Server) Handler(conn net.Conn) { ... }

// 启动服务器的接口
func (s *Server) Start() { ... }
```



## 网络基础：相关调用

> Go语言内置了一个强大的`net`包，用于处理网络相关的操作。主要有以下几个关键的组件：

### **net.Listener**

`net.Listener`接口用于监听传入的网络连接，通常用来创建服务器。例如，可以用`net.Listen`函数来监听特定的IP和端口。

```go
listener, err := net.Listen("tcp", ":8080")
```

这个接口包含以下方法：

- `Accept()`: 接受传入的连接。
- `Close()`: 关闭监听器。
- `Addr()`: 返回监听器的网络地址。

### **net.Conn**

`net.Conn`接口表示一个网络连接，既可以是服务器端接收到的连接，也可以是客户端发起的连接。它是全双工的，支持同时读写数据。

- `Read()`: 从连接中读取数据。
- `Write()`: 向连接写入数据。
- `Close()`: 关闭连接。
- `LocalAddr()`: 获取本地地址。
- `RemoteAddr()`: 获取远程地址。

```go
conn, err := net.Dial("tcp", "example.com:80")
```

### **net.Dial 和 net.Listen**

这两个函数用于创建连接和监听连接：

- `net.Dial`：用于客户端发起网络连接，支持TCP、UDP、HTTP等协议。
- `net.Listen`：用于服务器端监听连接。

### **net.TCPConn / net.UDPConn**

这两个是具体的网络协议连接对象，分别表示TCP和UDP连接。

- `net.TCPConn`: 用于TCP连接，继承自`net.Conn`。
- `net.UDPConn`: 用于UDP连接，支持无连接的数据报传输。
