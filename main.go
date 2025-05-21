package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ConnError struct {
	statusCode int
	message    string
	details    string
}

func (e *ConnError) Error() string {
	code := e.statusCode
	msg := e.message

	if code == 0 {
		code = 500
	}
	if msg == "" {
		msg = "Unknown error"
	}
	return fmt.Sprintf("STATUS: NOT OK\nSTATUS CODE: %d\nMESSAGE: %s\nDETAILS: %s", code, msg, e.details)
}

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		panic(err)
	}
	fmt.Println("Serving on PORT 8080")
	defer listener.Close()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Register for SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("Server is being closed due to signal: ", sig)
		listener.Close()
		done <- true
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			break
		}

		go func(c net.Conn) {
			defer c.Close()
			fmt.Println("Client connected: ", c.RemoteAddr())
			if connErr := handleConn(c); connErr == nil {
				_, err = c.Write([]byte("STATUS: OK\nMESSAGE: File received successfully\n"))

				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				_, err = c.Write([]byte(connErr.Error()))

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}(conn)
	}

	<-done
	fmt.Println("Server shutdown")
}

func handleConn(conn net.Conn) error {
	const MaxFileSize = 10 * 1024 * 1024 // 10 MB

	reader := bufio.NewReader(conn)
	filename := ""
	var filesize int = 0

	fmt.Println("START Reading incoming message...")
	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading message:", err.Error())
			return &ConnError{statusCode: 500, message: "Error in reading message", details: err.Error()}
		}

		line = strings.TrimRight(line, "\r\n")
		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)

		header := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		if len(parts) != 2 || (header != "filename" && header != "filesize") {
			fmt.Println("Malformed Headers")
			return &ConnError{statusCode: 400, message: "Bad Request", details: "Malformed Headers"}
		}

		switch header {
		case "filename":
			filename = value

		case "filesize":
			filesize, err = strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid filesize:", err)
				return &ConnError{statusCode: 400, message: "Bad Request", details: "Invalid filesize: " + err.Error()}
			}
		}

	}

	// filesize exceeding maximum file size
	if filesize > MaxFileSize {
		fmt.Println("Filesize too large")
		return &ConnError{statusCode: 400, message: "Bad Request", details: "File size exceeds " + strconv.Itoa(MaxFileSize) + " bytes"}
	}

	// cleaning the file path name
	filename = filepath.Base(filename)
	fmt.Printf("START Reading file content from %v of size %v\n", filename, filesize)

	// Generate unique prefix
	timestamp := time.Now().Unix()
	clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
	prefix := fmt.Sprintf("%d_%s", timestamp, clientIP)

	// Create the directory
	os.MkdirAll("goDropped", 0755)

	// Construct full path
	fullPath := filepath.Join("goDropped", fmt.Sprintf("%s_%s", prefix, filename))

	chunkSize := 4096 // 4KB
	received := 0
	buffer := make([]byte, chunkSize)

	file, err := os.Create(fullPath)
	if err != nil {
		return &ConnError{statusCode: 500, message: "Error creating file", details: err.Error()}
	}
	defer file.Close()

	for received < filesize {
		remaining := filesize - received
		if remaining < chunkSize {
			buffer = buffer[:remaining] // read only what's left
		}

		n, err := reader.Read(buffer)
		if err != nil {
			return &ConnError{statusCode: 500, message: "Error reading file content", details: err.Error()}
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			return &ConnError{statusCode: 500, message: "Error writing to file", details: err.Error()}
		}

		received += n

		// Simple progress bar
		percent := float64(received) / float64(filesize) * 100
		hashes := int(percent / 5)
		spaces := 20 - hashes
		fmt.Printf("\rReceiving: [%s%s] %.2f%%", strings.Repeat("#", hashes), strings.Repeat("-", spaces), percent)
	}

	fmt.Println("\nFile received successfully.")

	fmt.Printf("Saved %s in %s successfully\n", filename, fullPath)

	return nil
}
