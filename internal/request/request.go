package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	h "tcpToHttp/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
}

type Request struct {
	RequestLine RequestLine
	Headers     *h.Headers
	state       parserState
}

var ERROR_MALFORMED_REQ_LINE = fmt.Errorf("malformed request line.")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version.")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state.")
var CRLF = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: h.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]

		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}
			read += n

			if done {
				r.state = StateDone
			}
		case StateDone:
			break outer
		default:
			panic("someting went wrong in parse method")
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func parseRequestLine(r []byte) (*RequestLine, int, error) {
	idx := bytes.Index(r, CRLF)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := r[:idx]
	read := idx + len(CRLF)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQ_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   strings.TrimPrefix(string(parts[2]), "HTTP/"),
	}

	if !rl.ValidHTTP() {
		return nil, 0, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return rl, read, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// TODO : buffer size could get larger
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
