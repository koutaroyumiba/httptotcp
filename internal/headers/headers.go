package headers

import (
	"bytes"
	"fmt"
	"slices"
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

	if !isValidKey(key) {
		return read, done, fmt.Errorf("invalid key")
	}

	key = strings.ToLower(key)
	read += idx + len(CRLF)
	old_val, ok := h[key]
	if !ok {
		h[key] = value
	} else {
		h[key] = fmt.Sprintf("%s, %s", old_val, value)
	}

	return read, done, nil
}

func isValidKey(key string) bool {
	specialChar := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	if len(key) < 1 {
		return false
	}

	for _, c := range key {
		isAlpha := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
		isDigit := '0' <= c && c <= '9'
		isSpecialChar := slices.Contains(specialChar, c)

		if !isAlpha && !isDigit && !isSpecialChar {
			return false
		}
	}

	return true
}
