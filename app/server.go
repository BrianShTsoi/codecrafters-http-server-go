package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleGETRequest(conn net.Conn, request []string, dir string) {
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
		var b bytes.Buffer
		if os.IsNotExist(err) {
			b.WriteString("HTTP/1.1 404 Not Found\r\n\r\n")
		} else {
			b.WriteString("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\n")
			fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n%s", len(fileContent), fileContent)
		}
		conn.Write(b.Bytes())

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func handlePOSTRequest(conn net.Conn, request []string, dir string) {
	startLine := request[0]
	target := strings.Fields(startLine)[1]

	if target[0:7] == "/files/" {
		// fileContent := strings.Join(request[6:], "\r\n")
		fileContent := request[6]
		fileContentLen, _:= strconv.Atoi(request[3][16:])

		filepath := dir + "/" + target[7:]
		fp, _ := os.Create(filepath)
		defer fp.Close()
		fp.Write([]byte(fileContent)[:fileContentLen])

		var b bytes.Buffer
		b.WriteString("HTTP/1.1 201 Created\r\n\r\n")
		conn.Write(b.Bytes())
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()

	readBuf := make([]byte, 1024)
	_, err := conn.Read(readBuf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}

	request := strings.Split(string(readBuf), "\r\n")
	startLine := request[0]
	// target := strings.Fields(startLine)[1]

	if startLine[:3] == "GET" {
		handleGETRequest(conn, request, dir)
	} else if startLine[:4] == "POST" {
		handlePOSTRequest(conn, request, dir)
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
