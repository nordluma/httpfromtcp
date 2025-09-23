package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %s\n", err.Error())
	}
	defer file.Close()

	currLine := ""
	for {
		buf := make([]byte, 8)
		n, err := file.Read(buf)
		if err != nil {
			if currLine != "" {
				fmt.Printf("read: %s\n", currLine)
				currLine = ""
			}

			if err == io.EOF {
				break
			}

			fmt.Printf("error reading chunk: %s\n", err.Error())
			break
		}

		str := string(buf[:n])
		parts := strings.Split(str, "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", currLine+parts[i])
			currLine = ""
		}
		currLine += parts[len(parts)-1]

	}
}
