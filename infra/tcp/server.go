package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Server struct {
	listener net.Listener
	ch       chan []byte
}

func NewServer(address string, ch chan []byte) *Server {
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listen error: ", err)
		return nil
	}
	return &Server{listener: l, ch: ch}
}

func (server *Server) Serve() {
	fmt.Println("slog is listening on addr: ", server.listener.Addr())

	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			panic(err)
		}
		go receivePackets(conn, server.ch)
	}
}

func receivePackets(c net.Conn, ch chan []byte) {
	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadBytes('\000')
		if err != nil && err == io.EOF {
			c.Close()
			fmt.Printf("close connection [%v]\n", c)
			break
		} else if err != nil {
			panic(err)
		}
		ch <- msg
		ReceivePact++
	}
}

var ReceivePact int64