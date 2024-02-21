package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"gopkg.in/natefinch/npipe.v2"
)

func main() {
	ExampleListen()
}

// Use Listen to start a server, and accept connections with Accept().
func ExampleListen() {
	ln, err := npipe.Listen(`\\.\pipe\mypipe`)
	if err != nil {
		// handle error
		log.Fatal("something went wrong")
	}

	defer ln.Close()

	messageCount := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Println("some error with connection")
			continue
		}

		// handle connection like any other net.Conn
		go func(conn net.Conn) {
			defer conn.Close()
			defer fmt.Println()

			r := bufio.NewReader(conn)

			for {
				msg, err := r.ReadString('\n')
				if err != nil {
					// handle error or end of connection
					return
				}

				fmt.Print(msg)

				// If an empty string is received, emit an empty response
				if msg == "\n" {
					if _, err := fmt.Fprintln(conn, ""); err != nil {
						// handle error writing response
						return
					}

					// Close the connection after responding to the newline
					conn.Close()
					return
				} else {
					messageCount++
					response := fmt.Sprintf("I've received %d messages so far\n", messageCount)
					if _, err := fmt.Fprintln(conn, response); err != nil {
						// handle error writing response
						return
					}
				}
			}
		}(conn)
	}
}
