package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcpToHttp/internal/request"
	"tcpToHttp/internal/response"
	"tcpToHttp/internal/server"
)

const port = 42069

func main() {
	server := server.New(port)

	server.GET("/", defaultHandler)
	server.GET("/yourproblem", yourProblemHandler)
	server.GET("/myproblem", myProblemHandler)

	if err := server.Serve(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func defaultHandler(w io.Writer, req *request.Request) *server.HandlerError {
	w.Write([]byte("All good, frfr\n"))
	return nil
}

func yourProblemHandler(w io.Writer, req *request.Request) *server.HandlerError {
	return &server.HandlerError{
		StatusCode: response.StatusBadReq,
		Message:    "Your problem is not my problem\n",
	}
}

func myProblemHandler(w io.Writer, req *request.Request) *server.HandlerError {
	return &server.HandlerError{
		StatusCode: response.StatusServerError,
		Message:    "Woopsie, my bad\n",
	}
}
