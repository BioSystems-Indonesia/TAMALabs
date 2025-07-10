package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":1024")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 4096)
			for {
				n, err := c.Read(buf)
				if n > 0 {
					fmt.Print("Received message: ", string(buf[:n]))
				}
				if err != nil {
					if err.Error() != "EOF" {
						fmt.Println("Error reading:", err)
					}
					break
				}
			}
		}(conn)
	}
}
