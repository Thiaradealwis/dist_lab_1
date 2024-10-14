package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// all
	// Deal with an error event.
	if err != nil {
		fmt.Println("error")
	}

}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// So long as this connection is alive:
	reader := bufio.NewReader(client)
	for {
		// Read in new messages as delimited by '\n's
		msg, _ := reader.ReadString('\n')
		// Tidy up each message and add it to the messages channel,
		// recording which client it came from.
		msgs <- Message{clientid, msg}
	}

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	clientId := 0
	for {
		select {
		case conn := <-conns:

			// Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			clients[clientId] = conn
			// - start to asynchronously handle messages from this client
			handleClient(conn, clientId, msgs)
			clientId++
		case msg := <-msgs:
			// Deal with a new message
			// Send the message to all clients that aren't the sender
			msgs <- msg
		}
	}
}
