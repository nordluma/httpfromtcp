package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("error resolving udp address: %s\n", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("error connecting with udp: %s\n", err.Error())
		return
	}
	defer conn.Close()

	rdr := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := rdr.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading from stdin: %s\n", err.Error())
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("error writing to udp conn: %s\n", err.Error())
		}
	}
}
