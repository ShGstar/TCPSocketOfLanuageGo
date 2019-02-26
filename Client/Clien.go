package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	server := "localhost:7373"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)

	if err != nil {
		fmt.Println(os.Stderr, "Fatal error: ", err)
		os.Exit(1)
	}

	//建立服务器连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), os.Stderr, "Fatal error:", err)
		os.Exit(1)
	}

	fmt.Println("connection success")
	sender(conn)
	fmt.Println("send over")

}

func sender(conn *net.TCPConn) {
	words := "hello world!"
	msgBack, err := conn.Write([]byte(words)) //给服务器发信息

	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), "服务器反馈")
		os.Exit(1)
	}
	buffer := make([]byte, 1024)
	msg, err := conn.Read(buffer) //接受服务器信息
	fmt.Println(conn.RemoteAddr().String(), "服务器反馈：", string(buffer[:msg]), msgBack, "；实际发送了", len(words))
	conn.Write([]byte("ok")) //在告诉服务器，它的反馈收到了。
}

// func ResolveIPAddr
// func ResolveIPAddr(net, addr string) (*IPAddr, error)
// ResolveIPAddr parses addr as an IP address of the form "host" or "ipv6-host%zone"
// and resolves the domain name on the network net, which must be "ip", "ip4" or "ip6".
// ResolveIPAddr将addr解析为“host”或“ipv6-host％zone”形式的IP地址，并解析网络上的域名，
// 该域名必须是“ip”，“ip4”或“ip6”。

// func DialTCP
// func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)
// DialTCP connects to the remote address raddr on the network net,
// which must be "tcp", "tcp4", or "tcp6". If laddr is not nil, it is used as the local address for the connection.
