package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var CRLF = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	idx := bytes.Index(data, CRLF)
	if idx == -1 {
		return read, done, nil
	}

	if idx == 0 {
		done = true
		return read, done, nil
	}

	fieldPair := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(fieldPair) != 2 {
		return read, done, fmt.Errorf("invalid field format: %v", fieldPair)
	}

	key := string(fieldPair[0])
	value := strings.TrimSpace(string(fieldPair[1]))
	if len(key) != len(strings.TrimSpace(key)) {
		return read, done, fmt.Errorf("invalid field format (key spacing)")
	}

	read += idx + len(CRLF)
	h[key] = value
	return read, done, nil
}
