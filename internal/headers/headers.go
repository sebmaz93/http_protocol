package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("malformed field line.")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("malformed field name.")
var CRLF = []byte("\r\n")

func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if !found {
			return false
		}

	}
	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_FIELD_LINE
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(key, []byte(" ")) || !isToken(key) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}
	return string(key), string(value), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Set(key, value string) {
	lowerKey := strings.ToLower(key)
	oldValue := h.Get(lowerKey)
	if oldValue == "" {
		h.headers[lowerKey] = value
	} else {
		h.headers[lowerKey] = strings.Join([]string{oldValue, value}, ", ")
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], CRLF)
		if idx == -1 {
			break
		}

		// end of header
		if idx == 0 {
			done = true
			break
		}

		key, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(CRLF)
		h.Set(key, value)
	}
	return read, done, nil
}
