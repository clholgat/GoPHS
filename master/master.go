package main

import (
	"../ohhai"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
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
	request := &ohhai.OhHai{
		MessageType:      ohhai.OhHai_HEARTBEAT_REQUEST.Enum(),
		HeartBeatRequest: &ohhai.OhHai_HeartBeatRequest{},
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
