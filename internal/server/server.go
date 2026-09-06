package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/koutaroyumiba/httpfromtcp/internal/response"
)

type Server struct {
	Listener net.Listener
	closed   atomic.Bool
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: listener,
	}

	server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}

	_ = s.closed.Swap(true)
	return nil
}

func (s *Server) listen() {
	go func() {
		for {
			if s.closed.Load() {
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
	// res := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n")
	// _, err := conn.Write(res)
	// if err != nil {
	// 	log.Printf("error writing to the conn: %v", err)
	// }
	response.WriteStatusLine(conn, response.StatusOk)
	h := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, h)
}
