package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type serverState string

const (
	StateConnected    serverState = "connected"
	StateDisconnected serverState = "disconnected"
	StateError        serverState = "error"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port uint16) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: l,
	}

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"
	conn.Write([]byte(res))
	return
}
