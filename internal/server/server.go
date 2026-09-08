package server

import (
	"fmt"
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

func (h HandlerError) Write(w *response.Writer) error {
	if w.GetState() != response.Statusline {
		log.Printf("[HandlerErr] writer not in proper state: %s - %s", w.GetState(), h.Message)
		return fmt.Errorf("[HandlerErr] writer not in proper state: %s", w.GetState())
	}
	w.WriteStatusLine(h.Code)

	contentLen := len(h.Message)
	headers := response.GetDefaultHeaders(contentLen)
	w.WriteHeaders(headers)

	_, err := w.WriteBody([]byte(h.Message))
	return err
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: listener,
		Handler:  handler,
	}

	go server.listen()

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
}

func (s *Server) handle(conn net.Conn) {
	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("[ERROR] something went horrible: %v", err)
	}

	writer := response.NewResponseWriter(conn)
	handlerErr := s.Handler(writer, r)

	if handlerErr != nil {
		handlerErr.Write(writer)
		return
	}
}
