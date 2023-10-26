package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()

	readBuf := make([]byte, 1024)
	_, err := conn.Read(readBuf)
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
		conn.Write([]byte(writeBuf))

	} else if target == "/user-agent" {
		agent := request[2][12:]

		writeBuf := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(agent), agent)
		conn.Write([]byte(writeBuf))

	} else if target[0:7] == "/files/" {
		filepath := dir + "/" + target[7:]
		fileContent, err := os.ReadFile(filepath)
		// fmt.Printf("Request: %v\n", request)
		// fmt.Printf("Response: %v\n", writeBuf)
		var b bytes.Buffer
		if os.IsNotExist(err) {
			b.WriteString("HTTP/1.1 404 Not Found\r\n\r\n")

			fmt.Printf("Request: %v\n", request)
			fmt.Printf("filepath: %v\n", filepath)
			fmt.Printf("Response: %v\n", b.String())
		} else {
			b.WriteString("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\n")
			fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n%s", len(fileContent), fileContent)
			fmt.Printf("Request: %v\n", request)
			fmt.Printf("Response: %v\n", b.String())
		}


		conn.Write(b.Bytes())

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

func main() {
	dir := flag.String("directory", "", "")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, *dir)
	}
}
