package util

import (
	"../comm"
	"bufio"
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"io"
	"net"
)

func GetProto(conn net.Conn) (*comm.OhHai, error) {
	// Read the first 8 bytes.
	buf := make([]byte, 8)
	_, err := io.ReadAtLeast(conn, buf, 8)
	if err != nil {
		return nil, err
	}

	// Convert first 8 bytes to int
	lenbuf := bytes.NewBuffer(buf)
	var length int64
	err = binary.Read(lenbuf, binary.LittleEndian, &length)
	if err != nil {
		panic(err)
	}

	// Get the protobuf from the connection
	buf = make([]byte, length)
	_, err = io.ReadAtLeast(conn, buf, int(length))
	if err != nil {
		return nil, err
	}

	// Unmarshal the protobuf
	message := &comm.OhHai{}
	err = proto.Unmarshal(buf, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func SendProto(conn net.Conn, message *comm.OhHai) error {
	encoded, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(conn)

	// Convert the length of the marshalled protobuf to 
	// a byte array.
	size := int64(len(encoded))
	sizebytes := new(bytes.Buffer)
	err = binary.Write(sizebytes, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	// Send the size and then the protobuf.
	writer.Write(sizebytes.Bytes())
	writer.Write(encoded)
	writer.Flush()
	return nil
}
