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

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

type WriterState string

const (
	Statusline WriterState = "status-line"
	Headers    WriterState = "headers"
	Body       WriterState = "body"
)

type Writer struct {
	writer io.Writer
	state  WriterState
}

func NewResponseWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  Statusline,
	}
}

func (w Writer) GetState() WriterState {
	return w.state
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != Statusline {
		return fmt.Errorf("[WriteStatusLine] expected: %s, got: %s", Statusline, w.state)
	}
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
	_, err := w.writer.Write(statusLine)
	if err != nil {
		return err
	}

	w.state = Headers
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != Headers {
		return fmt.Errorf("[WriteStatusLine] expected: %s, got: %s", Headers, w.state)
	}
	resHeaders := []byte{}
	for key, value := range headers {
		resHeaders = fmt.Appendf(resHeaders, "%s: %s\r\n", key, value)
	}
	resHeaders = fmt.Append(resHeaders, "\r\n")

	_, err := w.writer.Write(resHeaders)
	if err != nil {
		return err
	}

	w.state = Body
	return nil
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	if w.state != Body {
		return 0, fmt.Errorf("[WriteStatusLine] expected: %s, got: %s", Body, w.state)
	}
	return w.writer.Write(b)
}
