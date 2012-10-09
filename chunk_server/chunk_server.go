package main

import (
	"../ohhai"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type Configuration struct {
	Master     string
	StorageDir string
}

var (
	SERVER_ID = 0
	config    *Configuration
	listener  net.Listener
)

// This function runs before main
// Handles setup with the master server to grab an ID
func init() {
	// Read master server from config file
	file, err := os.Open("server.cfg")
	if err != nil {
		panic(err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	config = &Configuration{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(config.StorageDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Creating file directory")
			err = os.Mkdir(config.StorageDir, os.ModeDir)
			os.Chmod(config.StorageDir, 0666)
			if err != nil {
				panic(err)
			}
		}
	}

	// Call master and ask for a server ID
	conn, err := net.Dial("tcp", config.Master)
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
	buf, err = ioutil.ReadAll(conn)
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
		heartBeatResponse(conn)
	case ohhai.OhHai_READ_REQUEST:
		fmt.Println("Read")
		readResponse(message.GetReadRequest(), conn)
	case ohhai.OhHai_WRITE_REQUEST:
		fmt.Println("Write")
	default:
		fmt.Println("WAT?")
	}
}

func heartBeatResponse(conn net.Conn) {
	fmt.Println("sending heartbeat response")
	heartbeat := &ohhai.OhHai_HeartBeatResponse{
		Id: make([]int64, 2, 10),
	}

	files, err := ioutil.ReadDir(config.StorageDir)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(files); i++ {
		file, err := strconv.Atoi(files[i].Name())
		if err == nil {
			heartbeat.Id = append(heartbeat.Id, int64(file))
		}
	}

	response := &ohhai.OhHai{
		MessageType:       ohhai.OhHai_HEARTBEAT_RESPONSE.Enum(),
		HeartBeatResponse: heartbeat,
	}

	writer := bufio.NewWriter(conn)
	bytes, err := proto.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, err = writer.Write(bytes)
	if err != nil {
		panic(err)
	}
	conn.Close()
}

func readResponse(request *OhHai_ReadRequest, conn net.Conn) {
	fmt.Println("sending read respnose")
	fileName := request.GetId()
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	upper := request.GetRangeTop()
	lower := request.GetRangeBottom()
	chunk := file[uper:lower]

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(chunk)
	if err != null {
		panic(err)
	}
	conn.Close()
}
