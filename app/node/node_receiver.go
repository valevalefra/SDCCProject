package main

import (
	"SDCCProject/app/utility"
	"log"
	"net"
)

func channel_for_message() {
	listener, err := net.Listen("tcp", ":"+"2345")
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {

	defer connection.Close()
	msg := new(utility.Message)

}
