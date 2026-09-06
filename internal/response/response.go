package response

import (
	"fmt"
	"io"

	"github.com/koutaroyumiba/httpfromtcp/internal/headers"
)

type StatusCode uint16

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := ""
	switch statusCode {
	case StatusOk:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "InternalServerError"
	}

	statusLine := fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, reason)
	_, err := w.Write(statusLine)
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	resHeaders := []byte{}
	for key, value := range headers {
		resHeaders = fmt.Appendf(resHeaders, "%s: %s\r\n", key, value)
	}
	resHeaders = fmt.Append(resHeaders, "\r\n")

	_, err := w.Write(resHeaders)
	if err != nil {
		return err
	}

	return nil
}
