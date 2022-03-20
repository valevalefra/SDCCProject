package main

import (
	"SDCCProject/app/utility"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func sendMessages(args ...string) error {
	var function func(msgs []string)
	function = send_to

	function(args)

	return nil

}

func send_to(msgs []string) {

	for _, text := range msgs {
		//increment local clock
		fmt.Printf("il nodo con id %d e valore del clock %d sta inviando %s \n", myId, *&scalarClock, text)
		incrementClock(&scalarClock)

		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.Request
		msg.Clock = getValueClock(&scalarClock)
		msg.Text = text
		msg.SendID = myId

		send_to_peer(msg)

	}
}

func send_to_peer(msg utility.Message) {

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utility.Info)
		//open connection whit other peer
		peer_conn := dest.Address + ":" + dest.Port
		conn, err := net.Dial("tcp", peer_conn)
		defer conn.Close()
		if err != nil {
			log.Println("Send response error on Dial")
		}
		enc := gob.NewEncoder(conn)
		enc.Encode(msg)
	}
}
