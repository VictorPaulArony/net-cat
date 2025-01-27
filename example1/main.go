package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// stract for the client properties
type Client struct {
	Conn     net.Conn
	UserName string
}

func main() {
	fmt.Println("Welcome to net-cat: ")

	// start the server for client calls
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Println("Error while listening to the port: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("Connection started at port 1234")
	fmt.Println("Clients can connect using: nc localhost 1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error while accepting the connection")
			continue
		}
		go HandleConnection(conn)
	}
}

// function to handle the client to server requests and calls response
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// enable client to enter their anme
	fmt.Fprintf(conn, "Enter User name: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	userName := scanner.Text()

	fmt.Fprintf(conn, "%s has joined the chat \n", userName)

	for scanner.Scan() {
		sms := scanner.Text()
		if sms == "quit" {
			fmt.Fprintf(conn, " %s has left the chat \n", userName)
			return
		}
		fmt.Printf("%s: %s \n", userName, sms)
		fmt.Fprintf(conn, "you: %s\n", sms)
	}
}
