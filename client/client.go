package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	// In a continuous loop, read a message from the server and display it.
	for {
		reader := bufio.NewReader(conn)
		msg, _ := reader.ReadString('\n')
		fmt.Printf(msg)
	}
}

func write(conn net.Conn) {
	// Continually get input from the user and send messages to the server.
	for {
		stdin := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter message")
		msg, _ := stdin.ReadString('\n')
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	// Try to connect to the server
	conn, _ := net.Dial("ip", *addrPtr)
	for {
		// Start asynchronously reading and displaying messages
		go read(conn)
		// Start getting and sending user messages.
		go write(conn)
	}
}
