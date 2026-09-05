tcp:
  go run ./cmd/tcplistener

udp:
  go run ./cmd/udpsender

test:
  go test ./...

ping-get:
  curl localhost:42069/i-use-neovim-btw

ping-post:
  curl -X POST localhost:42069/coffee \
      -H 'Content-Type: application/json' \
      -d '{"type": "dark mode", "size": "medium"}'
