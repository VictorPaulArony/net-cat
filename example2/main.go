package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Welcome to my net-cat program")

	if len(os.Args) != 2 {
		log.Fatalln("Invalid usage, use: go run [file name] [flag name]")
	}

	flag := os.Args[1]
	if flag == "-s" {
		Server()
	} else if flag == "-c" {
		Client()
	} else {
		fmt.Println("Invalid flag: use [-s] or [-c]")
		return
	}
}

// function to serve and listen to the connections and service request
func Server() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Println("Error while listening to the port", err)
		return
	}
	defer listener.Close()
	fmt.Println("Conection started at port :1234")

	// accepting the connection form the client when client is connected
	server, err := listener.Accept()
	if err != nil {
		fmt.Println("Failed to connect to the server: ", err)
	}
	defer server.Close()
	fmt.Println("Client accepted onnection started at port :1234")

	// starting a goroutine for the communication between the server and the client
	go func() {
		scanner := bufio.NewScanner(server)
		for scanner.Scan() {
			sms := scanner.Text()
			fmt.Println(sms)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sms := scanner.Text()
		fmt.Fprintf(server, "server: %s\n", sms)
	}
}

// function to connect the client to the server and call for service
func Client() {
	// client connecting to the server to call for the services
	client, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Println("Error while connecting to the server: ", err)
	}
	defer client.Close()
	fmt.Println("connected to the server successfully")

	// client requesting for services from the server
	go func() {
		scanner := bufio.NewScanner(client)
		for scanner.Scan() {
			sms := scanner.Text()
			fmt.Println(sms)
		}
	}()

	// client recive the service from the client call
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sms := scanner.Text()
		fmt.Fprintf(client, "Client: %s\n", sms)
	}
}
