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
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	str := string(buffer[:n])
	vecStr := strings.Split(str, "\r\n")
	path := strings.Split(vecStr[0], " ")[1]
	userAgent := ""
	if strings.HasPrefix(vecStr[1], "User-Agent") {
		userAgent = vecStr[1][len("User-Agent: "):]
	} else {
		userAgent = vecStr[2][len("User-Agent: "):]
	}
	if path == "/" {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			return
		}
	} else if strings.HasPrefix(path, "/echo/") {
		cont, _ := strings.CutPrefix(path, "/echo/")
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(cont), cont)
		_, err := conn.Write([]byte(resp))
		if err != nil {
			return
		}
	} else if strings.HasPrefix(path, "/user-agent/") {
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		_, err := conn.Write([]byte(resp))
		if err != nil {
			return
		}

	} else {
		_, err := conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			return
		}
	}
	err = conn.Close()
	if err != nil {
		return
	}
}
