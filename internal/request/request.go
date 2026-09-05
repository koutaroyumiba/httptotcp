package request

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	splitData := strings.Split(string(data), "\r\n")
	requestLine, err := parseRequestLine(splitData[0])
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return &Request{*requestLine}, nil
}

func parseRequestLine(requestLine string) (*RequestLine, error) {
	splitLine := strings.Split(requestLine, " ")
	if len(splitLine) != 3 {
		return nil, fmt.Errorf("length wrong")
	}

	method := splitLine[0]
	if !slices.Contains([]string{"GET", "POST", "DELETE", "PUT", "PATCH"}, method) {
		return nil, fmt.Errorf("bad method")
	}

	requestTarget := splitLine[1]

	version := strings.TrimPrefix(splitLine[2], "HTTP/")
	if version != "1.1" {
		return nil, fmt.Errorf("bad http version")
	}

	return &RequestLine{
		version,
		requestTarget,
		method,
	}, nil
}
