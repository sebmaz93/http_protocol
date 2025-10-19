package headers

import (
	"bytes"
	"fmt"
)

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("malformed field line.")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("malformed field name.")
var CRLF = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_FIELD_LINE
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}
	return string(key), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
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
		h[key] = value
	}
	return read, done, nil
}
