package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/koutaroyumiba/httpfromtcp/internal/request"
	"github.com/koutaroyumiba/httpfromtcp/internal/response"
)

type Server struct {
	Listener net.Listener
	Handler  Handler
	closed   bool
}

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: listener,
		Handler:  handler,
	}

	server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}

	s.closed = true
	return nil
}

func (s *Server) listen() {
	go func() {
		for {
			if s.closed {
				log.Printf("[warn] attempting to connect while server closed")
				return
			}

			conn, err := s.Listener.Accept()
			if err != nil {
				log.Printf("[error] conn err - %v", err)
				continue
			}

			s.handle(conn)

			conn.Close()
		}
	}()
}

func (s *Server) handle(conn net.Conn) {
	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("[ERROR] something went horrible: %v", err)
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerErr := s.Handler(buffer, r)
	if handlerErr != nil {
		WriteErrorResponse(conn, handlerErr)
		return
	}

	body := buffer.Bytes()
	response.WriteStatusLine(conn, response.StatusOk)
	h := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(conn, h)
	_, err = conn.Write(body)
	if err != nil {
		log.Fatalf("what happened here... %v", err)
	}
}

func WriteErrorResponse(w io.Writer, handlerErr *HandlerError) error {
	response.WriteStatusLine(w, handlerErr.Code)
	contentLen := len(handlerErr.Message)
	h := response.GetDefaultHeaders(contentLen)
	response.WriteHeaders(w, h)

	_, err := w.Write([]byte(handlerErr.Message))
	return err
}
