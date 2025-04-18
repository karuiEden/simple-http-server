package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var mapHttp = map[int32]string{
	200: "OK",
	201: "Created",
	404: "Not Found",
}

type Request struct {
	Method        string
	Path          string
	Version       string
	Host          string
	UserAgent     string
	Accept        string
	ContentType   string
	ContentLength int64
	Encoding      string
	Body          string
}

func newRequest(s string) Request {
	req := Request{}
	vecStr := strings.Split(s, "\r\n")
	requestLine := strings.Split(vecStr[0], " ")
	req.Method = requestLine[0]
	req.Path = requestLine[1]
	req.Version = requestLine[2]
	for _, str := range vecStr[1:] {
		if strings.HasPrefix(str, "Host") {
			req.Host = strings.Split(str, " ")[1]
		}
		if strings.HasPrefix(str, "User-Agent") {
			req.UserAgent = strings.Split(str, " ")[1]
		}
		if strings.HasPrefix(str, "Accept") {
			req.Accept = strings.Split(str, " ")[1]
		}
		if strings.HasPrefix(str, "Accept-Encoding") {
			req.Encoding = strings.Split(str, " ")[1]
		}
	}
	if len(vecStr[2]) != 0 {
		if strings.HasPrefix(vecStr[2], "Content-Type") {
			req.ContentType = strings.Split(vecStr[2], " ")[1]
		}
		if strings.HasPrefix(vecStr[2], "Content-Length") {
			req.ContentLength, _ = strconv.ParseInt(strings.Split(vecStr[2], " ")[1], 10, 64)
		}
	}
	req.Body = vecStr[len(vecStr)-1]
	return req
}

func buildResponse(code int32, req Request, contentType string, cont string) string {
	resp := fmt.Sprintf("%s %d %s\r\n", req.Version, code, mapHttp[code])
	if cont != "" {
		if req.Encoding == "gzip" {
			resp += fmt.Sprintf("Content-Encoding: gzip\r\n")
		}
		resp += fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n\r\n%s", contentType, len(cont), cont)
	} else {
		resp += "\r\n"
	}
	return resp
}

func rootHandler(r Request, conn net.Conn) error {
	var code int32
	if r.Path == "/" {
		code = 200
	} else {
		code = 404
	}
	resp := buildResponse(code, r, "", "")
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func userAgentHandler(r Request, conn net.Conn) error {
	resp := buildResponse(200, r, "text/plain", r.UserAgent)
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func echoHandler(r Request, conn net.Conn) error {
	cont, _ := strings.CutPrefix(r.Path, "/echo/")
	resp := buildResponse(200, r, "text/plain", cont)
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func fileHandler(r Request, conn net.Conn) error {
	filePath := os.Args[2] + strings.Split(r.Path, "/")[2]
	var resp string
	if r.Method == "GET" {
		cont, err := os.ReadFile(filePath)
		if err != nil {
			resp = buildResponse(404, r, "", "")
		} else {
			resp = buildResponse(200, r, "application/octet-stream", string(cont))
		}
	} else if r.Method == "POST" {
		err := os.WriteFile(filePath, []byte(r.Body), 0666)
		if err != nil {
			return err
		}
		resp = buildResponse(201, r, "", "")
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
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
