package utils

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// struct to store/track all the clients and there sms
type Sever struct {
	Clients    map[net.Conn]Client
	Messages   []Message
	Joined     chan Client
	Left       chan net.Conn
	Sms        chan Message
	MaxClients int
	Mu         sync.Mutex
}

// struct to define the client
type Client struct {
	Name string
	Conn net.Conn
}

// struct to define the Message content of the clients
type Message struct {
	Sender    net.Conn
	Sms       string
	Name      string
	Timestamp string
}

// function to start a new server for the client
func NewServer() *Sever {
	return &Sever{
		Clients:    make(map[net.Conn]Client),
		Messages:   []Message{},
		Joined:     make(chan Client),
		Left:       make(chan net.Conn),
		Sms:        make(chan Message),
		MaxClients: 10,
	}
}

// function to enable communication between the clients
func (s *Sever) HandleCommunication(conn net.Conn) {
	s.Mu.Lock()

	if len(s.Clients) >= s.MaxClients{
		fmt.Fprintln(conn, "Server currently full.Please try again later")
		conn.Close()
		return
	}
	s.Mu.Unlock()

	fmt.Fprintln(conn, "Welcome to TCP-Chat!\n"+
		"        _nnnn_\n"+
		"       dGGGGMMb\n"+
		"      @p~qp~~qMb\n"+
		"      M|@||@) M|\n"+
		"      @,----.JM|\n"+
		"     JS^\\__/  qKL\n"+
		"    dZP        qKRb\n"+
		"   dZP          qKKb\n"+
		"  fZP            SMMb\n"+
		"  HZM            MMMM\n"+
		"  FqM            MMMM\n"+
		" __| \".         |\\dS\"qML\n"+
		" |    `.        | `' \\Zq\n"+
		"_)       \\.___.,|     .'\n"+
		" \\____   )MMMMMP|   .'\n"+
		"      `-'       `--'")

	// promt new user connected to the system to ender there name
	fmt.Fprintf(conn, "[Enter User name]: ")

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	UserName := scanner.Text()

	s.Mu.Lock()
	for _, client := range s.Clients {
		if UserName == client.Name || UserName == "" {
			fmt.Fprintf(conn, "Username already taken or you cannot sign in with no username")
			conn.Close()
			s.Mu.Unlock()
			return
		}
	}

	client := Client{
		Name: UserName,
		Conn: conn,
	}
	s.Clients[conn] = client
	s.Mu.Unlock()


	s.Broadcast(client.Conn, "has joined our chat...")

	for _, prevMessages := range s.Messages {
		fmt.Fprintf(client.Conn, "[%s][%v]: %s \n", prevMessages.Timestamp, prevMessages.Name, prevMessages.Sms)
	}

	for scanner.Scan() {

		smss := scanner.Text()
		timer := time.Now().Format("2006-01-02 15:04:05")

		s.Mu.Lock()
		sms := Message{
			Sender:    client.Conn,
			Sms:       smss,
			Name:      UserName,
			Timestamp: timer,
		}

		s.Messages = append(s.Messages, sms)


		for _, client := range s.Clients {
			if client.Conn != sms.Sender {
				fmt.Fprintf(client.Conn, "[%s][%v]: %s \n", sms.Timestamp, sms.Name, sms.Sms)
			}
		}

		s.Mu.Unlock()

		if smss == "leave" {
			s.Broadcast(client.Conn, "has left our chat...")
			delete(s.Clients, conn)
			conn.Close()
			return
		}
	}

	// s.Left <- conn
	conn.Close()
}


func (s *Sever) Broadcast(sender net.Conn, sms string) {
	s.Mu.Lock()

	client := s.Clients[sender].Name
	s.Mu.Unlock()

	for conn := range s.Clients {
		if conn != sender { // this enables not to send sms back to sender
			fmt.Fprintf(conn, "%s %s \n", client, sms)
		}
	}
}
