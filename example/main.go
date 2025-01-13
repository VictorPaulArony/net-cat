package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// the client struct to store the user names and conns
type Client struct {
	UserName string
	Conn     net.Conn
	Address  string
}

// struct to manage all the clients connected
type Server struct {
	Clients    map[net.Conn]Client
	Mu         sync.Mutex
	MaxClients int
}

// function main to start the connections
func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	defaultPort := "8989"

	if len(os.Args) == 2 {
		defaultPort = os.Args[1]
	}

	// get the ip address of the server
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		// fmt.Println("Available on:")
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// if ipnet.IP.To4() != nil {
				// 	fmt.Printf("  http://%s:%s\n", ipnet.IP.String(), defaultPort)
				// }
			}
		}
	}
	listener, err := net.Listen("tcp", ":"+defaultPort)
	if err != nil {
		log.Printf("Failed to start server at port :%s: %v\n", defaultPort, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on the port :%s\n", defaultPort)

	server := CreatServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}
		go server.HandleConnection(conn)
	}
}

// function to create a new server connectioin for the clients
func CreatServer() *Server {
	return &Server{
		Clients:    make(map[net.Conn]Client),
		MaxClients: 10,
	}
}

// function to get the time for the clients sms
func GetTimestamp() string {
	return time.Now().Format("2006-10-10 15:14:05")
}

// method to displaye the chats of the clients sms
func (s *Server) Broadcast(sender net.Conn, sms string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	timestamp := GetTimestamp()
	senderName := s.Clients[sender].UserName

	for conn := range s.Clients {
		if conn != sender { // this enables not to send sms back to sender
			fmt.Fprintf(conn, "[%s][%s]: %s", timestamp, senderName, sms)
		}
	}
}

// method to handle the connection between multiple
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	s.Mu.Lock()
	if len(s.Clients) >= s.MaxClients {
		fmt.Fprintf(conn, "Server is full, maximum connections is %d \n", s.MaxClients)
		return
	}
	s.Mu.Unlock()

	// Enable user to login to the chart with unique user name
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
	fmt.Fprintf(conn, "[ENTER YOUR NAME]: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	userName := scanner.Text()

	// check if user name has whitespaces
	if userName == "" || strings.Contains(userName, " ") {
		fmt.Fprintln(conn, "Invalid username. Username cannot be empty or contain spaces.")
		return
	}

	// check if your name already exist in the chat room
	s.Mu.Lock()
	for _, client := range s.Clients {
		if client.UserName == userName {
			s.Mu.Unlock()
			fmt.Fprintln(conn, "Username already taken. Please choose another one.")
			return
		}
	}
	s.Mu.Unlock()

	clientAdrr := conn.RemoteAddr().String() // get client remote address

	// update new client when logged in
	client := Client{
		UserName: userName,
		Conn:     conn,
		Address:  clientAdrr,
	}

	// adding the client to the sever and to the db
	s.Mu.Lock()
	s.Clients[conn] = client
	s.Mu.Unlock()

	// broadcasting the client that has joined
	joinMessage := fmt.Sprintf("%s has joined our chat...", userName)
	s.Broadcast(conn, joinMessage)

	for scanner.Scan() {
		message := scanner.Text()

		if message == "" {
			continue
		}

		if message == "quit" {
			break
		}

		// Print the message locally for the sender
		timestamp := GetTimestamp()
		fmt.Fprintf(conn, "[%s][%s]:%s\n",
			timestamp,
			userName,
			message)

		// Broadcast to others (message won't be sent back to sender)
		s.Broadcast(conn, message)

	}

	s.Mu.Lock()
	delete(s.Clients, conn)
	s.Mu.Unlock()

	leaveMessage := fmt.Sprintf("%s has left our chat...", userName)
	s.Broadcast(conn, leaveMessage)
}
