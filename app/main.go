package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	Method    string
	Path      string
	Version   string
	Host      string
	UserAgent string
	Accept    string
	Body      string
}

func newRequest(s string) Request {
	req := Request{}
	vecStr := strings.Split(s, "\r\n")
	requestLine := strings.Split(vecStr[0], " ")
	req.Method = requestLine[0]
	req.Path = requestLine[1]
	req.Version = requestLine[2]
	for _, str := range vecStr {
		if strings.HasPrefix(str, "Host") {
			req.Host = strings.Split(str, " ")[1]
		}
		if strings.HasPrefix(str, "User-Agent") {
			req.UserAgent = strings.Split(str, " ")[1]
		}
		if strings.HasPrefix(str, "Accept") {
			req.Accept = strings.Split(str, " ")[1]
		}
	}
	req.Body = vecStr[len(vecStr)-1]
	return req
}

func rootHandler(r Request, conn net.Conn) error {
	var resp string
	if r.Path == "/" {
		resp = r.Version + "200 OK" + "\r\n\r\n"
	} else {
		resp = r.Version + "404 Not Found" + "\r\n\r\n"
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func userAgentHandler(r Request, conn net.Conn) error {
	resp := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", r.Version, len(r.UserAgent), r.UserAgent)
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func echoHandler(r Request, conn net.Conn) error {
	cont, _ := strings.CutPrefix(r.Path, "/echo/")
	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(cont), cont)
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func fileHandler(r Request, conn net.Conn) error {
	if r.Method == "GET" {
		filePath := os.Args[2] + strings.Split(r.Path, "/")[2]
		cont, err := os.ReadFile(filePath)
		var resp string
		if err != nil {
			resp = fmt.Sprintf("%s 404 Not Found\r\n\r\n", r.Version)
		} else {
			resp = fmt.Sprintf("%s 200 OK\r\nContent-Type: octet-stream\r\nContent-Length: %d\r\n\r\n%s", r.Version, len(cont), cont)
		}
		_, err = conn.Write([]byte(resp))
		if err != nil {
			return err
		}
	}
	return nil
}

func Handler(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	str := string(buffer[:n])
	req := newRequest(str)
	if strings.HasPrefix(req.Path, "/echo/") {
		err := echoHandler(req, conn)
		if err != nil {
			return
		}
	} else if strings.HasPrefix(req.Path, "/user-agent") {
		err := userAgentHandler(req, conn)
		if err != nil {
			return
		}
	} else if strings.HasPrefix(req.Path, "/files") && len(os.Args) > 2 && os.Args[1] == "--directory" {
		err := fileHandler(req, conn)
		if err != nil {
			return
		}
	} else {
		err := rootHandler(req, conn)
		if err != nil {
			return
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage

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
		go Handler(conn)

	}

}
