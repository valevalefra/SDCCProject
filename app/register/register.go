package main

import (
	"SDCCProject/app/utility"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

func main() {
	var connectNum int
	reg := new(utility.Utility)

	server := rpc.NewServer()
	//register method
	err := server.RegisterName("Register", reg)
	if err != nil {
		log.Fatal("Format of service Register is not correct: ", err)
	}

	port := 4321
	// listen for a connection
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	defer listener.Close()

	log.Printf("RPC server on port %d", port)

	go server.Accept(listener)

	//Wait connection
	for connectNum < 3 {
		ch := <-utility.Connection
		if ch == true {
			connectNum++
		}
	}

	log.Printf("Max Number of Connection reached up")

	utility.Wg.Add(-3)

	select {}
}
