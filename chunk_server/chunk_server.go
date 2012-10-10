package main

import (
	"../comm"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	//"io"
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

	// Send a come alive to master
	address := listener.Addr().String()
	alive := &comm.OhHai{
		MessageType: comm.OhHai_COME_ALIVE.Enum(),
		ComeAlive:   &comm.ComeAlive{Server: &address},
	}
	bytes, err := proto.Marshal(alive)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(conn)
	fmt.Println("connecting to master")
	// Send the listening address to master
	writer.Write(bytes)
	writer.WriteByte(0)
	writer.Flush()

	reader := bufio.NewReader(conn)
	fmt.Println(len(buf))
	buf, err = reader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
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
	message := &comm.OhHai{}
	err = proto.Unmarshal(buf, message)
	message_type := message.GetMessageType()
	switch message_type {
	case comm.OhHai_HEARTBEAT_REQUEST:
		fmt.Println("Hearbeat")
		heartBeatResponse(conn)
	case comm.OhHai_READ_REQUEST:
		fmt.Println("Read")
		readResponse(message.GetReadRequest(), conn)
	case comm.OhHai_WRITE_REQUEST:
		fmt.Println("Write")
	default:
		fmt.Println("WAT?")
	}
}

func heartBeatResponse(conn net.Conn) {
	fmt.Println("sending heartbeat response")
	heartbeat := &comm.HeartBeatResponse{
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

	response := &comm.OhHai{
		MessageType:       comm.OhHai_HEARTBEAT_RESPONSE.Enum(),
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

func readResponse(request *comm.ReadRequest, conn net.Conn) {
	fmt.Println("sending read respnose")
	fileName := request.GetId()
	file, err := os.Open(string(fileName))
	if err != nil {
		panic(err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	upper := request.GetRangeTop()
	lower := request.GetRangeBottom()
	chunk := buf[upper:lower]

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(chunk)
	if err != nil {
		panic(err)
	}
	conn.Close()
}
