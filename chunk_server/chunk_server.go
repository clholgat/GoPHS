package main

import (
	"github.com/clholgat/GoPHS/ohhai"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"net"
	"os"
)

var (
	SERVER_ID = 0
	MASTER    = ""
)

// This function runs before main
// Handles setup with the master server to grab an ID
func init() {
	// Read master server from config file
	file, err := os.Open("master.cfg")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	MASTER, err := reader.ReadString('\n')
	if err != nil {
		//panic(err)
	}
	fmt.Println(MASTER)

	// Call master and ask for a server ID
	conn, err := net.Dial("tcp", MASTER)
	if err != nil {
		panic(err)
	}
	fmt.Println("sending ping to master")
	fmt.Fprintf(conn, "ping")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}

func main() {
	fmt.Println("All the things")

	l, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	buf, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	message := &ohhai.OhHai{}
	err = proto.Unmarshal(buf, message)
	message_type := message.GetMessageType()
	switch message_type {
	case ohhai.OhHai_HEARTBEAT_REQUEST:
		fmt.Println("Hearbeat")
	case ohhai.OhHai_READ_REQUEST:
		fmt.Println("Read")
	case ohhai.OhHai_WRITE_REQUEST:
		fmt.Println("Write")
	default:
		fmt.Println("WAT?")
	}
}

func getChunkId() int64 {
	return int64(1)
}
