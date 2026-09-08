package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/koutaroyumiba/httpfromtcp/internal/headers"
)

var CRLF = []byte("\r\n")

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

	switch r.state {
	case StateInit:
		parsedRequestLine, n, err := parseRequestLine(data)
		if err != nil {
			log.Printf("error: %v", err)
			return read, err
		}

		if n == 0 {
			break
		}

		read += n
		r.RequestLine = *parsedRequestLine
		r.state = StateHeaders

	case StateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			log.Printf("[header err] %v", err)
			return read, err
		}

		if done {
			r.state = StateBody
		}

		read += n

	case StateBody:
		lengthStr, ok := r.Headers["content-length"]
		if !ok {
			r.state = StateDone
			break
		}

		contentLength, err := strconv.Atoi(lengthStr)
		if err != nil {
			log.Printf("error: %v", err)
			return read, err
		}

		if contentLength == 0 {
			r.state = StateDone
			break
		}

		dataLength := len(data)
		if len(r.Body)+dataLength > contentLength {
			log.Printf("error: length of body greater than Content-Length\n")
			return read, fmt.Errorf("length of body greater than Content-Length")
		}

		r.Body = slices.Concat(r.Body, data)
		if len(r.Body) == contentLength {
			r.state = StateDone
			break
		}

		read += dataLength

	case StateDone:
		break
	}

	return read, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		Headers: headers.NewHeaders(),
		state:   StateInit,
	}

	buffer := make([]byte, 4096)
	read := 0

	for request.state != StateDone {
		n, err := reader.Read(buffer[read:])
		if err != nil {
			log.Printf("(RequestFromReader) [error] %v\n", err)
			return nil, err
		}

		read += n
		processed := 1
		for processed != 0 {
			processed, err = request.parse(buffer[:read])
			if err != nil {
				log.Printf("[parseErr] %v", err)
				return nil, err
			}

			copy(buffer, buffer[processed:read])
			read -= processed
		}
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, CRLF)
	if idx == -1 {
		return nil, 0, nil
	}

	// process
	requestLine := string(data[:idx])
	read := idx + len(CRLF)

	splitLine := strings.Split(requestLine, " ")
	if len(splitLine) != 3 {
		return nil, read, fmt.Errorf("wrong length %s", splitLine)
	}

	method := splitLine[0]
	if !slices.Contains([]string{"GET", "POST", "DELETE", "PUT", "PATCH"}, method) {
		return nil, read, fmt.Errorf("bad method (got %s)", method)
	}

	requestTarget := splitLine[1]

	version := strings.TrimPrefix(splitLine[2], "HTTP/")
	if version != "1.1" {
		return nil, read, fmt.Errorf("bad http version (got %s)", version)
	}

	return &RequestLine{
		version,
		requestTarget,
		method,
	}, read, nil
}
