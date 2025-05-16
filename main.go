package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		panic(err)
	}
	fmt.Println("Serving on PORT 8080")
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			fmt.Println("Client connected: ", c.RemoteAddr())
			if handleConn(c) {
				_, err = c.Write([]byte("STATUS: OK\nMESSAGE: File received successfully\n"))

				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				_, err = c.Write([]byte("STATUS: NOT OK\nMESSAGE: Error in receiving file\n"))

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}(conn)
	}

}

func handleConn(conn net.Conn) bool {

	reader := bufio.NewReader(conn)
	filename := ""
	var filesize int = 0

	_ = filename
	_ = filesize

	readFileContent := false

	fmt.Println("START Reading incoming message...")
	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading message:", err.Error())
			return false
		}

		line = strings.TrimSpace(line)

		if line == "" && !readFileContent {
			readFileContent = true
			break
		}

		if !readFileContent {
			parts := strings.SplitN(line, ":", 2)

			if len(parts) != 2 {
				fmt.Println("Malformed Headers")
				return false
			}

			header := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			switch header {
			case "filename":
				filename = value

			case "filesize":
				filesize, err = strconv.Atoi(value)
				if err != nil {
					fmt.Println("Invalid filesize:", err)
					return false
				}
			}
		} else {
			break
		}
	}

	fmt.Printf("START Reading file content from %v of size %v\n", filename, filesize)

	// code for reading file content
	buf := make([]byte, filesize)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println("Error in receiving file")
		return false
	}

	fmt.Println("File content: ", string(buf[:n]))

	return true
}
