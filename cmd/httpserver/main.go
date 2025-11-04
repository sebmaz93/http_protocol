package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
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

func defaultHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	body := []byte(`
		<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`)
	headers := response.GetDefaultHeaders(0)
	headers.Set("Content-length", strconv.Itoa(len(body)), true)
	headers.Set("Content-type", "text/html", true)
	res.WriteStatusLine(response.StatusOK)
	res.WriteHeaders(*headers)
	res.WriteBody(body)
	return nil
}

func yourProblemHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	body := []byte(`
		<html>
			<head>
    			<title>400 Bad Request</title>
      		</head>
        	<body>
         		<h1>Bad Request</h1>
           		<p>Your request honestly kinda sucked.</p>
            </body>
        </html>
	`)
	headers := response.GetDefaultHeaders(0)
	headers.Set("Content-length", strconv.Itoa(len(body)), true)
	headers.Set("Content-type", "text/html", true)
	res.WriteStatusLine(response.StatusBadReq)
	res.WriteHeaders(*headers)
	res.WriteBody(body)
	return &server.HandlerError{
		StatusCode: response.StatusBadReq,
		Message:    "Your problem is not my problem\n",
	}
}

func myProblemHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	body := []byte(`
		<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)

	headers := response.GetDefaultHeaders(0)
	headers.Set("Content-length", strconv.Itoa(len(body)), true)
	headers.Set("Content-type", "text/html", true)
	res.WriteStatusLine(response.StatusServerError)
	res.WriteHeaders(*headers)
	res.WriteBody(body)
	return &server.HandlerError{
		StatusCode: response.StatusServerError,
		Message:    "Woopsie, my bad\n",
	}
}
