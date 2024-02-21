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
	conn, err := npipe.Dial(`\\.\pipe\mypipe`)
	if err != nil {
		log.Fatal("Error connecting to the named pipe:", err)
	}

	message := fmt.Sprintf("Hi server! Count: %d", count)

	// Create a channel to receive responses from the goroutine
	responseChan := make(chan string)

	// Goroutine for reading and buffering responses
	go func() {
		for {
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Fatal("Error reading response from the named pipe:", err)
			}

			// Check for an empty response and close the connection
			if response == "\n" {
				conn.Close()
				close(responseChan)
				return
			}

			// Send response to the channel
			responseChan <- response

		}
	}()

	// Send messages to the server
	for i := 0; i < count; i++ {
		if _, err := fmt.Fprintln(conn, message); err != nil {
			log.Fatal("Error writing to the named pipe:", err)
		}
	}

	// send newline to the server
	if _, err := fmt.Fprintln(conn); err != nil {
		log.Fatal("Error writing to the named pipe:", err)
	}

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
