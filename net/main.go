package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/victorpaularony/net/utils"
)

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Unable to listen to that port: ", err)
		return
	}

	defer listener.Close()
	fmt.Printf("Listening on the port :%s\n", port)

	server := utils.NewServer()

	// recursion to allow the communication of the clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection at that port: ", err)
			continue
		}
		go server.HandleCommunication(conn)
	}
}
