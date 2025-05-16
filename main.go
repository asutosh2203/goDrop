package main

import (
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {

		n, err := conn.Read(buf)

		if err == nil {
			fmt.Println(string(buf[:n]))
		}

		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello"))

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
