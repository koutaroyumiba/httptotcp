package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/koutaroyumiba/httpfromtcp/internal/request"
	"github.com/koutaroyumiba/httpfromtcp/internal/response"
	"github.com/koutaroyumiba/httpfromtcp/internal/server"
)

const port = 42069

func get400() string {
	return `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
}

func get500() string {
	return `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
}

func get200() string {
	return `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
}

func customHandler(w *response.Writer, r *request.Request) *server.HandlerError {
	responseStatus := response.StatusOk
	responseHeaders := response.GetDefaultHeaders(0)
	responseBody := get200()

	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		responseStatus = response.StatusBadRequest
		responseBody = get400()
	case "/myproblem":
		responseStatus = response.StatusInternalServerError
		responseBody = get500()
	}

	w.WriteStatusLine(responseStatus)
	responseHeaders.Replace("content-type", "text/html")
	responseHeaders.Replace("content-length", string(len(responseBody)))
	w.WriteHeaders(responseHeaders)
	w.WriteBody([]byte(responseBody))

	return nil
}

func main() {
	server, err := server.Serve(port, customHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
