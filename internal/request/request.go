package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_MALFORMED_REQ_LINE = fmt.Errorf("malformed request line.")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version.")
var SEPARATOR = "\r\n"

func parseRequestLine(r string) (*RequestLine, string, error) {
	idx := strings.Index(r, SEPARATOR)
	if idx == -1 {
		return nil, r, nil
	}

	startLine := r[:idx]
	restOfMsg := r[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_REQ_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}

	if !rl.ValidHTTP() {
		return nil, restOfMsg, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return rl, restOfMsg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to io.ReadAll"),
			err,
		)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}
