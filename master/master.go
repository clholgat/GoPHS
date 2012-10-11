package main

import (
	"../comm"
	"../util"
	"fmt"
	//"io/ioutil"
	"net"
	"runtime"
)

var (
	CHUNK_SERVERS = make([]string, 0, 10)
)

func main() {
	go talk()
	listen()
}

func listen() {
	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Accepting")
		conn, err := l.Accept()
		fmt.Println("Accepted")

		if err != nil {
			panic(err)
		}

		go func(c net.Conn) {
			fmt.Println("received connection")

			message, err := util.GetProto(c)
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
				alive := message.GetComeAlive()
				server := alive.GetServer()
				fmt.Println(server, "is alive!")
				CHUNK_SERVERS = append(CHUNK_SERVERS, server)
				sendAck(c)
				sendHeartBeatRequest(server)
			default:
				fmt.Println("WAT?")
			}

			c.Close()
			fmt.Println("end connection")
		}(conn)
	}
	fmt.Println("done listening")
}

func talk() {
	for {
		for i := 0; i < len(CHUNK_SERVERS); i++ {
			sendHeartBeatRequest(CHUNK_SERVERS[i])
		}
		runtime.Gosched()
	}
}

func sendAck(conn net.Conn) {
	ack := &comm.OhHai{
		MessageType: comm.OhHai_ACK.Enum(),
	}

	err := util.SendProto(conn, ack)
	if err != nil {
		fmt.Println(err)
	}
}

// Send a request for a heartbeat to a given server.
func sendHeartBeatRequest(server string) {
	fmt.Println("sending heartbeat request")
	request := &comm.OhHai{
		MessageType:      comm.OhHai_HEARTBEAT_REQUEST.Enum(),
		HeartBeatRequest: &comm.HeartBeatRequest{},
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("could not contact", server, err)
		removeServer(server)
		return
	}

	err = util.SendProto(conn, request)
	if err != nil {
		fmt.Println("could not contact", server, err)
		removeServer(server)
		return
	}
	response, err := util.GetProto(conn)
	if err != nil {
		fmt.Println("could not contact", server, err)
		removeServer(server)
		return
	}
	fmt.Println("recieved heartbeat response")
	if response.GetMessageType() != *comm.OhHai_HEARTBEAT_RESPONSE.Enum() {
		panic("wrong response from heartbeat request")
	}
	chunks := response.GetHeartBeatResponse().Id
	for i := 0; i < len(chunks); i++ {
		fmt.Println(chunks[i])
	}
	conn.Close()
}

func removeServer(server string) {
	serverList := make([]string, len(CHUNK_SERVERS)-1)
	j := 0
	for i := 0; i < len(CHUNK_SERVERS); i++ {
		if server != CHUNK_SERVERS[i] {
			serverList[j] = CHUNK_SERVERS[i]
			j += 1
		}
	}
	CHUNK_SERVERS = serverList
}
