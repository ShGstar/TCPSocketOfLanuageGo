package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func ErrorOperator(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//客户端
func main() {

	conn, err := net.Dial("tcp4", "localhost:8000")
	ErrorOperator(err)
	defer func() {
		conn.Close()
	}()

	//创建一个协程输入，并发送消息
	go handleSend(conn)

	//在一个独立的协程中接受服务端的消息

	go handleReceive(conn)

	chanQuit := make(chan bool, 0)

	<-chanQuit
}

func handleReceive(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer) //没有消息会阻塞
		if err != io.EOF {
			ErrorOperator(err)
		}

		if n > 0 {
			msg := string(buffer[:n])
			fmt.Println(msg)
		}

	}
}

func handleSend(conn net.Conn) {
	//TODO:发送昵称到服务端
	//为了方便 此处利用随机数建立用户
	rand.Seed(time.Now().UnixNano())

	name := "无名氏" + strconv.Itoa(rand.Intn(100))

	_, err := conn.Write([]byte(name))
	ErrorOperator(err)

	input := bufio.NewReader(os.Stdin)

	for {
		line, _, _ := input.ReadLine()

		_, err := conn.Write(line)

		ErrorOperator(err)

		//正常退出
		if string(line) == "exit" {
			os.Exit(0)
		}
	}

}
