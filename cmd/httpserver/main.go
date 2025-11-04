package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	h "tcpToHttp/internal/headers"
	"tcpToHttp/internal/request"
	"tcpToHttp/internal/response"
	"tcpToHttp/internal/server"
)

const port = 42069

func toStr(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}
	return out
}

func main() {
	server := server.New(port)

	server.GET("/", defaultHandler)
	server.GET("/yourproblem", yourProblemHandler)
	server.GET("/myproblem", myProblemHandler)
	server.GET("/httpbin/stream/:count", chunkHandler)
	server.GET("/httpbin/:type", tailersHandler)
	server.GET("/video", videoHandler)

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

func chunkHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	parts := strings.Split(req.RequestLine.RequestTarget, "/")
	if len(parts) < 4 {
		return &server.HandlerError{
			StatusCode: response.StatusBadReq,
			Message:    "no count \n",
		}
	}

	_, err := strconv.Atoi(parts[3])
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusServerError,
			Message:    "Woopsie, my bad\n",
		}
	}

	binRes, err := http.Get("https://httpbin.org/stream/" + parts[3])
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	headers := response.GetDefaultHeaders(0)
	headers.Delete("Content-length")
	headers.Set("Transfer-Encoding", "chunked", true)
	headers.Set("Trailer", "X-Content-SHA256", false)
	headers.Set("Trailer", "X-Content-Length", false)
	res.WriteStatusLine(response.StatusOK)
	res.WriteHeaders(*headers)

	fullBody := []byte{}
	for {
		data := make([]byte, 32)
		n, err := binRes.Body.Read(data)
		if err != nil {
			break
		}
		// TODO: find more effiecent way
		fullBody = append(fullBody, data[:n]...)
		res.WriteBody([]byte(fmt.Sprintf("%x%s", n, request.CRLF)))
		res.WriteBody(data[:n])
		res.WriteBody(request.CRLF)
	}
	res.WriteBody([]byte(fmt.Sprintf("0%s", request.CRLF)))
	return nil
}

func tailersHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	parts := strings.Split(req.RequestLine.RequestTarget, "/")
	if len(parts) < 3 {
		return &server.HandlerError{
			StatusCode: response.StatusBadReq,
			Message:    "no html \n",
		}
	}

	binRes, err := http.Get("https://httpbin.org/" + parts[2])
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	headers := response.GetDefaultHeaders(0)
	headers.Delete("Content-length")
	headers.Set("Transfer-Encoding", "chunked", true)
	headers.Set("Trailer", "X-Content-SHA256", false)
	headers.Set("Trailer", "X-Content-Length", false)
	res.WriteStatusLine(response.StatusOK)
	res.WriteHeaders(*headers)

	fullBody := []byte{}
	for {
		data := make([]byte, 32)
		n, err := binRes.Body.Read(data)
		if err != nil {
			break
		}
		// TODO: find more effiecent way
		fullBody = append(fullBody, data[:n]...)
		res.WriteBody([]byte(fmt.Sprintf("%x%s", n, request.CRLF)))
		res.WriteBody(data[:n])
		res.WriteBody(request.CRLF)
	}
	res.WriteBody([]byte(fmt.Sprintf("0%s", request.CRLF)))
	trailers := h.NewHeaders()
	sha256.Sum256(fullBody)
	// TODO : remove hard coded sha and length the service was temp not available while did the testing
	trailers.Set("X-Content-SHA256", "3f324f9914742e62cf082861ba03b207282dba781c3349bee9d7c1b5ef8e0bfe", false)
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", 3741), false)
	res.WriteHeaders(*trailers)
	return nil
}

func videoHandler(res *response.Writer, req *request.Request) *server.HandlerError {
	file, err := os.Open("assets/vim.mp4")
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusServerError,
			Message:    "error reading video\n",
		}
	}
	defer file.Close()

	headers := response.GetDefaultHeaders(0)
	headers.Delete("Content-length")
	headers.Set("Transfer-Encoding", "chunked", true)
	headers.Set("Content-type", "video/mp4", true)
	headers.Set("Accept-Ranges", "none", true)
	res.WriteStatusLine(response.StatusOK)
	res.WriteHeaders(*headers)

	buffer := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			res.WriteBody([]byte(fmt.Sprintf("%x%s", n, request.CRLF)))
			res.WriteBody(buffer[:n])
			res.WriteBody(request.CRLF)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading video file: %v", err)
			break
		}
	}

	res.WriteBody([]byte(fmt.Sprintf("0%s", request.CRLF)))
	res.WriteBody(request.CRLF)
	return nil
}
