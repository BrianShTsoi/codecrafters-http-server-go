package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	readBuf := make([]byte, 1024)
	_, err = conn.Read(readBuf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}

	request := strings.Split(string(readBuf), "\r\n")
	startLine := request[0]
	target := strings.Fields(startLine)[1]


	if target == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	} else if target[0:6] == "/echo/" {
		random_str := target[6:]

		writeBuf := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(random_str), random_str)
		fmt.Printf("Request: %v\n", request)
		fmt.Printf("Response: %v\n", writeBuf)
		conn.Write([]byte(writeBuf))

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
