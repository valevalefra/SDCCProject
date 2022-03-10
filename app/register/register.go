package main

import (
	"SDCCProject/app/utils"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

type Register struct{}

func main() {
	var connect_num int
	utility := new(utils.Utils)

	server := rpc.NewServer()
	//register method
	err := server.RegisterName("Register", utility)
	if err != nil {
		log.Fatal("Format of service Utility is not correct: ", err)
	}

	port := 4321
	log.Println("ciaoooo44444")
	// listen for a connection
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	// Close the listener whenever we stop
	defer listener.Close()

	log.Printf("RPC server on port %d", port)

	go server.Accept(listener)

	//Wait connection
	for connect_num < 3 {
		ch := <-utils.Connection
		if ch == true {
			connect_num++
		}
	}

	log.Printf("Max Number of Connection reached up")

	utils.Wg.Add(-3)
	//send client a responce for max number of peer registered

	for {
		//TODO after registration this peer must be off ??
	}
}
