package main

import (
	"../comm"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
	//"io/ioutil"
	"net"
)

var (
	CHUNK_SERVERS = make([]string, 0, 10)
)

func main() {
	listen()
	talk()
}

func listen() {
	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func(c net.Conn) {
			fmt.Println("received connection")
			reader := bufio.NewReader(c)
			buf, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			fmt.Println("Read bytes", len(buf))

			// Strip the newline off the buffer
			//buf = buf[0 : len(buf)-1]
			fmt.Println(buf[len(buf)])
			message := &comm.OhHai{}
			err = proto.Unmarshal([]byte(buf), message)
			if err != nil {
				panic(err)
			}
			message_type := message.GetMessageType()
			switch message_type {
			case comm.OhHai_READ_REQUEST:
				fmt.Println("read request")
			case comm.OhHai_WRITE_REQUEST:
				fmt.Println("write reqest")
			case comm.OhHai_COME_ALIVE:
				fmt.Println("come alive")
			default:
				fmt.Println("WAT?")
			}

			server, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				panic(err)
			}
			server = server[0 : len(server)-1]
			fmt.Println("connecting. . .", server)
			CHUNK_SERVERS = append(CHUNK_SERVERS, server)
			sendHeartBeatRequest(server)
			io.WriteString(c, "1")
			c.Close()
		}(conn)
	}
}

func talk() {

}

func sendHeartBeatRequest(server string) {
	fmt.Println("sending heartbeat request")
	request := &comm.OhHai{
		MessageType:      comm.OhHai_HEARTBEAT_REQUEST.Enum(),
		HeartBeatRequest: &comm.HeartBeatRequest{},
	}
	fmt.Println(server)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(conn)
	bytes, err := proto.Marshal(request)
	if err != nil {
		panic(err)
	}
	_, err = writer.Write(bytes)
	if err != nil {
		panic(err)
	}
	conn.Close()
}
