package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		var currentLine strings.Builder

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Fatalf("Failed to read from file: %v", err)
				break
			}

			parts := strings.Split(string(data[:n]), "\n")
			if len(parts) == 2 {
				// new line found
				out <- fmt.Sprintf("%s%s", currentLine.String(), parts[0])
				currentLine.Reset()
				currentLine.WriteString(parts[1])
			} else if len(parts) == 1 {
				currentLine.WriteString(parts[0])
			} else {
				log.Fatal("something is wrong?")
			}
		}
		if currentLine.Len() != 0 {
			out <- fmt.Sprintf("%s", currentLine.String())
		}
	}()

	return out
}
