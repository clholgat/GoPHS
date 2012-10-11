package main

import (
	"../comm"
	"../util"
	"bufio"
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
	// Grab the config file.
	file, err := os.Open("server.cfg")
	if err != nil {
		panic(err)
	}

	// Read the config file.
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Parse the config file.
	config = &Configuration{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		panic(err)
	}

	// Check if the file directory exists.
	_, err = os.Stat(config.StorageDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Create directory if it doesn't exist.
			fmt.Println("Creating file directory")
			err = os.Mkdir(config.StorageDir, os.ModeDir)
			os.Chmod(config.StorageDir, 0666)
			if err != nil {
				panic(err)
			}
		}
	}

	// Call master for initial communication.
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

	// Build a come alive message.
	address := listener.Addr().String()
	alive := &comm.OhHai{
		MessageType: comm.OhHai_COME_ALIVE.Enum(),
		ComeAlive:   &comm.ComeAlive{Server: &address},
	}

	// Send the proto to master
	util.SendProto(conn, alive)

	// Get the response from the server.
	response, err := util.GetProto(conn)
	if err != nil {
		panic(err)
	}
	message_type := response.GetMessageType()
	switch message_type {
	case comm.OhHai_ACK:
		fmt.Println("Successfully connected to master")
	case comm.OhHai_ERROR:
		fmt.Println("Could not connect")
	default:
		fmt.Println("THESE ARE NOT THE DROIDS YOU'RE LOOKING FOR")
	}

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
	fmt.Println("new connection")

	message, err := util.GetProto(conn)
	if err != nil {
		panic(err)
	}
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
		Id: make([]int64, 0, 10),
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

	util.SendProto(conn, response)
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
