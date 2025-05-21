# 📁 GoDrop - Simple TCP File Receiver in Go

GoDrop is a lightweight TCP server written in Go that receives files over a custom TCP protocol and stores them in a local directory. Perfect for learning low-level networking and file transfer mechanics.

## 🚀 Features

- Accepts incoming TCP connections on `127.0.0.1:8080`
- Handles multiple clients concurrently using goroutines
- Supports basic header-based protocol (`filename` and `filesize`)
- **Displays a real-time progress bar** while receiving files
- Graceful shutdown on SIGINT / SIGTERM (e.g., Ctrl+C)
- Automatically stores received files in the `goDropped/` folder with a timestamped, IP-prefixed filename
- Max file size limit of **10 MB**

## 📦 Protocol Format

Client should send headers followed by the file bytes:

```
filename: my_file.txt
filesize: 12345

[RAW FILE BYTES...]
```

- Headers must be followed by a blank line (`\n`) before the file data begins.
- `filename` is the original file name.
- `filesize` is the number of bytes to be read from the stream after the headers.

## 📊 Progress Bar in Action

While receiving a file, GoDrop now shows a visual progress bar in the terminal:

`Receiving: [##########----------] 50.00%`

This updates in real time as the server writes chunks to disk.

## 🛠️ Running the Server

Make sure you have Go installed. Then run:

```bash
go run main.go
```
Server will start listening on `127.0.0.1:8080`

You should see:

```bash
Serving on PORT 8080
START Reading incoming message...
START Reading file content from my_file.txt of size 12345
Receiving: [####################] 100.00%
File received successfully.
Saved my_file.txt in goDropped/1716271874_127.0.0.1_my_file.txt successfully
```

## 🧪 Example Client (Using netcat)

You can test the server using nc (netcat) as a basic client:

```bash
(
  echo -e "filename: test.txt\nfilesize: $(wc -c < test.txt)\n\n";
  cat test.txt
) | nc 127.0.0.1 8080
```

Make sure `test.txt` exists in the same directory.

## 📁 Output

Received files will be stored in the `goDropped/` directory with names like:

```
goDropped/1716271874_127.0.0.1_test.txt
```

This includes a Unix timestamp and sender's IP to avoid collisions.

## 🧹 Graceful Shutdown

The server listens for system signals like `SIGINT` and `SIGTERM` to shut down cleanly:

```arduino
^C
Server is being closed due to signal: interrupt
Stopping accepting connections...
Server shutdown
```

## 🧯 Error Handling

- Malformed headers? You'll get a 400 Bad Request.
- File too big? The server politely declines.
- Disk I/O errors? They're logged, and the connection is safely dropped.

## 🧰 Tech Stack

- 🐹 Go (Golang)
- 📡 TCP Sockets
- 🧵 Goroutines
- 💾 os, bufio, net, and time packages

## 📜 License

MIT License. Do whatever you want. Just don't claim you wrote it first 😉.

---

Enjoy dropping files the Go way! 🚀
