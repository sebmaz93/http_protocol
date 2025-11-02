package response

import (
	"fmt"
	"io"
	"strconv"
	"tcpToHttp/internal/headers"
)

type StatusCode uint16

const (
	StatusOK               StatusCode = 200
	StatusBadReq           StatusCode = 400
	StatusNotFound         StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405
	StatusServerError      StatusCode = 500
)

var statusMap = map[StatusCode]string{
	StatusOK:               "OK",
	StatusBadReq:           "Bad Request",
	StatusNotFound:         "Not Found",
	StatusMethodNotAllowed: "Method not allowed",
	StatusServerError:      "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusText, ok := statusMap[statusCode]
	if !ok {
		statusText = "Unknown Status"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, header headers.Headers) error {
	b := []byte{}
	header.ForEach(func(k, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", k, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.Write(b)
	return err
}
