package main

import (
	"fmt"
	"log"
)

var (
	listener GameConnection
)

func main() {
	fmt.Println("Enter the IP of the remote host...")
	var host string
	fmt.Scanln(&host)
	fmt.Println("Enter the port of the remote host...")
	var port string
	fmt.Scanln(&port)
	remoteHost = host
	remotePort = port
	initConnections()

}

func initConnections() {
	listener = GameConnection{}
	log.Println("Starting listener...")
	InitListener(&listener)
}
