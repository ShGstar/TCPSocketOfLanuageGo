package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

//聊天
//服务器端

type ClientData struct {
	Conn *net.TCPConn
	Name string
	Addr string
}

func ErrorOperator(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//客户端信息，用昵称为键
//clientsMap = make(map[string]net.Conn)
var clientsMap = make(map[string]ClientData)

func main() {

	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8000")
	ErrorOperator(err)

	//建立监听
	lister, err := net.ListenTCP("tcp4", addr)
	ErrorOperator(err)

	defer func() {
		lister.Close()
	}()

	for { //循环接入所有
		conn, err := lister.AcceptTCP()
		ErrorOperator(err)

		clientAddr := conn.RemoteAddr()

		//TODO:接受并且保存昵称等
		buffer := make([]byte, 1024)

		//客户端名字
		var clientName string

		for { //循环读入数据
			n, err := conn.Read(buffer)
			ErrorOperator(err)
			//每个用户第一次发送的是名字等信息
			if n > 0 {
				clientName = string(buffer[:n])
				break
			}
		}
		fmt.Println(clientName + "上线了")

		//TODO:将每一个丢入MAP中
		Newclient := ClientData{conn, clientName, clientAddr.String()}
		clientsMap[clientName] = Newclient

		//TODO:给已经在线的用户发送上限通知

		//单独协程中聊天
		go ChatWithClient(Newclient)
	}

}

//和一个client 通信
func ChatWithClient(client ClientData) {

	buffer := make([]byte, 1024)

	for {
		n, err := client.Conn.Read(buffer)
		if err != io.EOF {
			ErrorOperator(err)
		}

		if n > 0 {
			msg := string(buffer[:n])
			fmt.Printf("%s:%s\n", client.Name, msg)

			//将客户端说的每一句话记录在log中
			WriteMsgToLog(msg, client)

			//标志
			str := strings.Split(msg, "#")
			if len(str) > 1 {
				//all#......

				//获取目标
				targetName := str[0]
				targetMsg := str[1]

				//TODO：群发还是单发
				if targetName == "all" {
					//群发消息
					for _, c := range clientsMap {
						c.Conn.Write([]byte(client.Name + ":" + targetMsg))
					}
				} else {
					//单发
					for key, c := range clientsMap {
						if key == targetName {
							c.Conn.Write([]byte(client.Name + ":" + targetMsg))

							//写入日志

							break
						}
						//没有key

					}

				}

			} else {
				//客户端主动下线
				if msg == "exit" {
					//将当前客户端从在线用户中除名
					//向其他用户发送下线通知
					for name, c := range clientsMap {
						if c == client {
							delete(clientsMap, name)
						} else {
							c.Conn.Write([]byte(name + "已经下线"))
						}
					}

				} else if strings.Index(msg, "log@") == 0 {
					//log@all
					//log@XXX
					//filterName := strings.Split(msg, "@")[1]
					filename := strings.Split(msg, "@")[1]

					sendLog2Client(client, filename)
					//向客户端发送聊天日志
				} else {
					client.Conn.Write([]byte("已阅:" + msg))
				}
			}

		}

	}

}

//对话记录到文件中
func WriteMsgToLog(msg string, client ClientData) {
	//判断文件是否存在
	_, err := os.Stat(client.Name + ".log")
	//var file *os.File

	if err != nil {
		fmt.Println("文件不存在，创建该用户的文件")
		file, errC := os.Create(client.Name + ".log")
		ErrorOperator(errC)
		defer file.Close()

		logMsg := fmt.Sprintln(time.Now().Format("2016-01-01 00:00:00"), msg)
		file.Write([]byte(logMsg))

	} else {
		//文件已经存在
		file, errO := os.OpenFile(client.Name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

		ErrorOperator(errO)
		defer file.Close()

		logMsg := fmt.Sprintln(time.Now().Format("2016-01-01 12:34:56"), msg, " ")
		_, err := file.Write([]byte(logMsg))
		ErrorOperator(err)
		//file.WriteString(logMsg)
	}
}

//向客户端发送它的聊天记录
func sendLog2Client(client ClientData, filterName string) error {
	//1，发送它的所有聊天记录
	//2,发送它于某人的所有聊天记录

	//判断文件是否存在
	_, err := os.Stat(client.Name + ".log")
	if err != nil {
		fmt.Println("文件不存在")
		return err
	}

	//目前都是发送所有
	msg, err := ioutil.ReadFile(client.Name + ".log")
	ErrorOperator(err)
	if filterName != "all" {
		client.Conn.Write(msg)
	} else {
		//发送所有
		client.Conn.Write(msg)
	}

	return nil
}
