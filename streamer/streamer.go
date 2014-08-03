package streamer

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type Stream struct {
	Path string
}

func NewStreamer(path string) *Stream {
	return &Stream{Path: path}
}

func (stream *Stream) Stream(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	fmt.Printf("Serving on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:\n\t", err)
			continue
		}
		fmt.Printf("Serving to %s\n", conn.RemoteAddr().String())

		go stream.Serve(conn)
	}
}

func (stream *Stream) Serve(conn net.Conn) {
	file, err := os.Open(stream.Path)
	if err != nil {
		fmt.Println("Error reading file:\n\t", err)
		return
	}

	bytes := make([]byte, 1024)

	conn.Write([]byte("HTTP/1.1 200 OK\n"))
	conn.Write([]byte("Content-Type: audio/ogg\n"))
	conn.Write([]byte("Connection: close\n"))
	conn.Write([]byte("\n"))

	var position int64 = 0
	for numRead, err := file.ReadAt(bytes, position); err != io.EOF; numRead, err = file.ReadAt(bytes, position) {
		position += int64(numRead)

		numWritten, err := conn.Write(bytes)
		if err != nil {
			fmt.Println("Error writing data:\n\t", err)
			conn.Close()
			break
		}
		fmt.Printf("Wrote %d bytes to connection.\n", numWritten)

		time.Sleep(10 * time.Millisecond)
	}

	conn.Close()
}
