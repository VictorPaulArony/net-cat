package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// type structs to help define the program
type Client struct {
	Conn     net.Conn
	UserName string
}

type Server struct {
	Clients    map[net.Conn]Client
	Mu         sync.Mutex
	MaxClients int
}

func main() {
	fmt.Println("Welcome to my Net-cat")

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Println("Failed to start server at port :1234 ", err)
		return
	}
	defer listener.Close()

	server := StartNewServer()
	fmt.Println("Chat started at port :1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}
		go server.HandleConnection(conn)
	}
}

// function to set/ start a new server for the connections
func StartNewServer() *Server {
	return &Server{
		Clients:    make(map[net.Conn]Client),
		MaxClients: 10,
	}
}

// method to enable brodcastng of the sms to all the clients
func (s *Server) Broadcast(sender net.Conn, sms string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for conn := range s.Clients {
		if conn != sender { // this enables not to send sms back to sender
			fmt.Fprintf(conn, "%s: %s\n", s.Clients[sender].UserName, sms)
		}
	}
}

// method to handle connection between the clients
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// first check if the connection is at max
	s.Mu.Lock()
	if len(s.Clients) >= s.MaxClients {
		s.Mu.Unlock()
		fmt.Fprintf(conn, "server is full, maximum connection is %d", s.MaxClients)
		return
	}

	// enabling the user to entr there user names and add them to the server
	fmt.Fprint(conn, "Enter Username: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	userName := scanner.Text()

	client := Client{
		Conn:     conn,
		UserName: userName,
	}
	s.Clients[conn] = client
	s.Mu.Unlock()

	// broadcasting the new user to all the users
	s.Broadcast(conn, "has joined the chat")

	// handling the clients sms
	for scanner.Scan() {
		sms := scanner.Text()
		if sms == "quit" {
			break
		}
		s.Broadcast(conn, sms)
	}

	// Remove client and announce departure
	s.Mu.Lock()
	delete(s.Clients, conn)
	s.Mu.Unlock()
	s.Broadcast(conn, "has left the chat")
}
