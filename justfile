http:
  go run ./cmd/httpserver

tcp:
  go run ./cmd/tcplistener

udp:
  go run ./cmd/udpsender

udp-client:
  nc -u -l 42069

test:
  go test ./...

ping-get:
  curl -v localhost:42069/i-use-neovim-btw

ping-post:
  curl -v -X POST localhost:42069/coffee \
      -H 'Content-Type: application/json' \
      -d '{"type": "dark mode", "size": "medium"}'

ping-400:
  curl -v localhost:42069/yourproblem

ping-500:
  curl -v localhost:42069/myproblem
