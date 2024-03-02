package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/natefinch/npipe.v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./your_executable <count>")
		os.Exit(1)
	}

	countStr := os.Args[1]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Fatal("Invalid count value:", err)
	}

	ExampleDial(count)
}

func ExampleDial(count int) {
	//TODO: name the pipe appropriately, and add timeout.
	conn, err := npipe.Dial(`\\.\pipe\mypipe`)
	if err != nil {
		log.Fatal("Error connecting to the named pipe:", err)
	}

	// define the buffer reader and writer
	buffWriter := bufio.NewWriter(conn)
	buffReader := bufio.NewReader(conn)

	// Create a channel to receive responses from the goroutine
	responseChan := make(chan string)

	// write messages to server
	writeMessagesToBuffer(count, buffWriter, fmt.Sprintf("Hi server! Count: %d", count))

	// Flush the buffer to ensure all messages are sent
	if err := buffWriter.Flush(); err != nil {
		_ = conn.Close()
		log.Fatal("Error flushing the buffer:", err)
	}

	// Goroutine for reading and buffering responses
	go handleResponses(buffReader, conn, responseChan)

	// Receive and print the last response
	// Create a variable to keep track of the maximum channel size
	counter := 0
	lastResponse := ""
	for response := range responseChan {
		lastResponse = response
		counter++
	}

	// Print the last response
	fmt.Printf("%slength: %d\n", lastResponse, counter)
}

// TODO: add a timeout to the response channel
func handleResponses(buffReader *bufio.Reader, conn *npipe.PipeConn, outputChan chan string) {
	for {
		response, err := buffReader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading response from the named pipe:", err)
		}

		// Check for an empty response and close the connection
		if response == "\n" {
			_ = conn.Close()
			close(outputChan)
			return
		} else {
			// Print the response received from the server
			fmt.Print(response)
		}

		// Send response to the channel
		outputChan <- response

	}
}

func writeMessagesToBuffer(count int, buffWriter *bufio.Writer, messageToServer string) {
	// Send messages to the server
	for i := 0; i < count; i++ {
		if _, err := fmt.Fprintln(buffWriter, messageToServer); err != nil {
			log.Fatal("Error writing to the named pipe:", err)
		}
	}
	// send newline to the server
	if _, err := fmt.Fprintln(buffWriter); err != nil {
		log.Fatal("Error writing to the named pipe:", err)
	}
}
