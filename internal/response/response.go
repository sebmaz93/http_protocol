package response

import (
	"fmt"
	"io"
	"strconv"
	"tcpToHttp/internal/headers"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusText, ok := statusMap[statusCode]
	if !ok {
		statusText = "Unknown Status"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	_, err := w.writer.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen), false)
	h.Set("Connection", "close", false)
	h.Set("Content-Type", "text/plain", false)
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{}
	headers.ForEach(func(k, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", k, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, err
}

// func (w *Writer) WriteChunkedBody(p []byte) (int, error)

// func (w *Writer) WriteChunkedBodyDone() (int, error)
