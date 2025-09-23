package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
	}
	defer file.Close()

	for {
		buf := make([]byte, 8)
		if _, err = file.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
		}

		fmt.Printf("read: %s\n", buf)
	}
}
