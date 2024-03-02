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
	log.Println("Starting server...")
	ln, err := npipe.Listen(`\\.\pipe\mypipe`)
	if err != nil {
		// handle error
		log.Fatal("something went wrong", err.Error())
	}
	log.Println("Server started!")

	defer ln.Close()

	messageCount := 0

	for {
		conn, err := ln.Accept()
		log.Println("Connection accepted!")
		if err != nil {
			// handle error
			log.Println("some error with connection")
			continue
		}

		// handle connection like any other net.Conn
		go func(conn net.Conn) {
			defer log.Println("Connection closed!")
			defer conn.Close()

			reader := bufio.NewReader(conn)
			writer := bufio.NewWriter(conn)

			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					// handle error or end of connection
					log.Println("connection ended or read fail occurred.")
					return
				}

				// Print the message received from the client
				fmt.Print(msg)

				if handleMessage(msg, writer, &messageCount, err) {
					return
				}
			}
		}(conn)
	}
}

func handleMessage(msg string, writer *bufio.Writer, messageCount *int, err error) bool {
	// If an empty string is received, emit an empty response
	if msg == "\n" {
		if _, err := fmt.Fprintln(writer, fmt.Sprintf("Server: I've received %d messages so far\n", *messageCount)); err != nil {
			// handle error writing response
			log.Println("Error writing response to the named pipe:", err.Error())
		}
		if _, err := fmt.Fprintln(writer, ""); err != nil {
			// handle error writing response
			log.Println("Error writing newline response to the named pipe:", err.Error())
		}
		err := writer.Flush()
		if err != nil {
			log.Println("Error while flushing")
		}
		return true // break out of the loop
	} else {
		*messageCount = *messageCount + 1

		if err != nil {
			log.Println("Error writing to the named pipe:", err.Error())
		}
	}
	return false
}
