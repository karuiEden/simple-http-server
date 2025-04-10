package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage

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
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return
	}
	str := string(buffer)
	vecStr := strings.Split(str, "\r\n")
	path := strings.Split(vecStr[0], " ")[1]
	if path == "/" || path == "/index.html" {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\n\r"))
		if err != nil {
			return
		}
	} else {
		_, err := conn.Write([]byte("HTTP/1.1 404 Not Found\n\r"))
		if err != nil {
			return
		}
	}
	err = conn.Close()
	if err != nil {
		return
	}
}
