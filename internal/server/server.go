package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"tcpToHttp/internal/request"
	"tcpToHttp/internal/response"
)

type HandlerFunc func(w *response.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (hErr HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(hErr.StatusCode)
	messageBytes := []byte(hErr.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	w.WriteHeaders(*headers)
	w.WriteBody(messageBytes)
}

type routeKey struct {
	method string
	path   string
}

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	port     uint16
	routes   map[routeKey]HandlerFunc
	mu       sync.RWMutex
}

func New(port uint16) *Server {
	return &Server{
		port:   port,
		routes: make(map[routeKey]HandlerFunc),
	}
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.registerRoute("GET", path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.registerRoute("POST", path, handler)
}

func (s *Server) registerRoute(method, path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := routeKey{
		method,
		path,
	}
	s.routes[key] = handler
	log.Printf("Registered %s %s", method, path)

}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	s.listener = l

	log.Printf("Server listening on port %d", s.port)
	go s.listen()
	return nil
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
		go s.handleConn(connection)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	resWriter := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadReq,
			Message:    err.Error(),
		}
		hErr.Write(resWriter)
		return
	}

	s.mu.RLock()
	key := routeKey{
		method: req.RequestLine.Method,
		path:   req.RequestLine.RequestTarget,
	}
	handler, exists := s.routes[key]

	if !exists {
		handler, exists = s.findPatternMatch(req.RequestLine.Method, req.RequestLine.RequestTarget)
	}
	s.mu.RUnlock()

	if !exists {
		s.mu.RLock()
		pathExists := false

		for k := range s.routes {
			if k.path == req.RequestLine.RequestTarget {
				pathExists = true
				break
			}
		}
		s.mu.RUnlock()

		if pathExists {
			hErr := &HandlerError{
				StatusCode: response.StatusMethodNotAllowed,
				Message:    err.Error(),
			}
			hErr.Write(resWriter)
		} else {
			hErr := &HandlerError{
				StatusCode: response.StatusNotFound,
				Message:    err.Error(),
			}
			hErr.Write(resWriter)
		}
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := handler(resWriter, req)
	if hErr != nil {
		hErr.Write(resWriter)
		return
	}

	b := buf.Bytes()
	headers := response.GetDefaultHeaders(len(b))
	resWriter.WriteStatusLine(response.StatusOK)
	resWriter.WriteHeaders(*headers)
	conn.Write(b)
	return
}

func (s *Server) findPatternMatch(method, path string) (HandlerFunc, bool) {
	for k, h := range s.routes {
		if k.method != method {
			continue
		}

		if strings.Contains(k.path, ":") {
			if s.matchPath(k.path, path) {
				return h, true
			}
		}
	}
	return nil, false
}

func (s *Server) matchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			continue
		}
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}
