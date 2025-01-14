package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// struncts to store the client and server info
type Client struct {
	UserName string
	Conn     net.Conn
}

type Server struct {
	Clients        map[net.Conn]Client
	JoinCh         chan Client
	SmsCh          chan Message
	LeaveCh        chan net.Conn
	MaxClients     int
	Mu             sync.Mutex
	MessageHistory []Message // Store history of messages
}

type Message struct {
	Sender net.Conn
	Sms    string
}

func main() {
	// this enables the server to listen for calls
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("Error while listening from port: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("connection started on port: 8080")

	// define the server to start the server goroutine for the connections
	server := StartNewServer()
	go server.Start()

	// handling the goroutine from the accepted connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error while accepting the connection: ", err)
			continue
		}
		go server.HandleConnection(conn)
	}
}

// function to start the server connection for the clients
func StartNewServer() *Server {
	return &Server{
		Clients:        make(map[net.Conn]Client),
		JoinCh:         make(chan Client),
		SmsCh:          make(chan Message),
		LeaveCh:        make(chan net.Conn),
		MaxClients:     10,
		MessageHistory: []Message{},
	}
}

// function to enable the client to communicate
func (s *Server) HandleConnection(conn net.Conn) {
	// enable the user to enter the connsction with user name
	fmt.Fprintf(conn, "Enter User Name: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	userName := scanner.Text()

	// uupdating the registered user and notify the rest
	client := Client{
		UserName: userName,
		Conn:     conn,
	}

	// Notify the server to add the new client to the connection
	s.JoinCh <- client

	// Send all previous messages to the new client
	s.Mu.Lock()
	for _, sms := range s.MessageHistory {
		fmt.Fprintf(conn, "%s: %s\n", s.Clients[sms.Sender].UserName, sms.Sms)
	}
	s.Mu.Unlock()

	// handling of the client sms after joining the connection
	for scanner.Scan() {
		sms := scanner.Text()
		if sms == "quit" {
			break
		}
		s.BroadcastSms(conn, sms)
	}

	// notify the server that the client is leaving the connection
	s.LeaveCh <- conn
}

// function to enable broadcasting of the sms to all the clients and not the sender
func (s *Server) BroadcastSms(sender net.Conn, sms string) {
	s.SmsCh <- Message{
		Sender: sender,
		Sms:    sms,
	}
}

// method to enable broadcasting of sms to all the clients
func (s *Server) Broadcast() {
	for sms := range s.SmsCh {
		s.Mu.Lock()

		// Store the message in history
		s.MessageHistory = append(s.MessageHistory, sms)

		for conn := range s.Clients {
			if conn != sms.Sender {
				fmt.Fprintf(conn, "[%s]: %s\n", s.Clients[sms.Sender].UserName, sms.Sms)
			}
		}
		s.Mu.Unlock()
	}
}

// function method to handle thr new clients joining and leaving the connection
func (s *Server) ManageClients() {
	for {
		select {
		case client := <-s.JoinCh:
			s.Mu.Lock()
			if len(s.Clients) >= s.MaxClients {
				s.Mu.Unlock()
				client.Conn.Write([]byte("Server is full. Please try again later.\n"))
				client.Conn.Close()
				continue
			}
			s.Clients[client.Conn] = client
			s.Mu.Unlock()
			s.BroadcastSms(client.Conn, "has joined the chat")

		case conn := <-s.LeaveCh:
			s.Mu.Lock()
			delete(s.Clients, conn)
			s.Mu.Unlock()
			s.BroadcastSms(conn, "has left the chat")
		}
	}
}

// function method to start the server to handle all goruotines
func (s *Server) Start() {
	go s.Broadcast()
	go s.ManageClients()
	select {}
}
