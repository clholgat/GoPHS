package main

import (
	"../ohhai"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

var (
	SERVER_ID = 0
	MASTER    = ""
	listener  net.Listener
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

	// Listen for communication from master on an arbitrary
	// open port.
	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	fmt.Println("sending ping to master")
	// Send the listening address to master
	io.WriteString(conn, listener.Addr().String()+"\n")
	buf, err := ioutil.ReadAll(conn)
	status := string(buf)
	fmt.Println(status)
	conn.Close()
}

func main() {
	fmt.Println("Running main chunk server")

	for {
		// Wait for messages
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}

}

// Handle all of the connections.
func handleConnection(conn net.Conn) {
	buf, err := ioutil.ReadAll(conn)
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
