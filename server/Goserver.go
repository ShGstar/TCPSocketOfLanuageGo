package main

import "net"
import "fmt"

func main() {
	//服务器
	listener, err := net.Listen("tcp", "localhost:7373")
	if err != nil {
		fmt.Println("Listen is err!：", err)
	}

	for {
		conn, err := listener.Accept() //开启监听
		if err != nil {
			fmt.Println("Accept is err!: ", err)
			continue
		}
		//发生了连接
		fmt.Println("tcp connect success:", conn.RemoteAddr().String())
		//go handleConnection(conn)
		handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		//接受客户端消息
		msg, err := conn.Read(buffer)

		if err != nil {
			//接受错误
			fmt.Println("connection err!:", err)
			return
		}
		//接受正确
		fmt.Print(conn.RemoteAddr().String())
		fmt.Println("receive data: ", string(buffer[:msg]))

		//反馈给客户端
		bufferReturn := "我收到了"
		msgW, errW := conn.Write([]byte(bufferReturn))

		//确认客户端没有收到回执
		if errW != nil {
			fmt.Print(conn.RemoteAddr().String(), msgW)
			fmt.Println("没有收到回执")
			return
		}

		//确认客户端收到回执
		msg, err = conn.Read(buffer)
		fmt.Println(conn.RemoteAddr().String(), "客户端收到回执", string(buffer[:msg]), "客户收到了", msgW, "；实际发送了", len(bufferReturn))
	}
	defer conn.Close()
}

//type Conn interface {
// Read从连接中读取数据
// Read方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
//Read(b []byte) (n int, err error)

// Write从连接中写入数据
// Write方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
//Write(b []byte) (n int, err error)  func (c *TCPConn) Write(b []byte) (int, error)

// Close方法关闭该连接
// 并会导致任何阻塞中的Read或Write方法不再阻塞并返回错误
//Close() error

// 返回本地网络地址
//LocalAddr() Addr

// 返回远端网络地址
//RemoteAddr() Addr

// 设定该连接的读写deadline，等价于同时调用SetReadDeadline和SetWriteDeadline
// deadline是一个绝对时间，超过该时间后I/O操作就会直接因超时失败返回而不会阻塞
// deadline对之后的所有I/O操作都起效，而不仅仅是下一次的读或写操作
// 参数t为零值表示不设置期限 SetDeadline(t time.Time) error
// 设定该连接的读操作deadline，参数t为零值表示不设置期限
//SetReadDeadline(t time.Time) error

// 设定该连接的写操作deadline，参数t为零值表示不设置期限
// 即使写入超时，返回值n也可能>0，说明成功写入了部分数据
//SetWriteDeadline(t time.Time) error
//}

//type Listener interface {
// Accept waits for and returns the next connection to the listener.
// Accept() (c Conn, err error)

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
// Close() error

// Addr returns the listener's network address.
// Addr() Addr
//}

// func Listen
// func Listen(net, laddr string) (Listener, error)
// Listen announces on the local network address laddr. The network net must be a stream-oriented network:
// "tcp", "tcp4", "tcp6", "unix" or "unixpacket". See Dial for the syntax of laddr.
