package main

import (
	"SDCCProject/app/utility"
	"encoding/gob"
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
		incrementClock(&scalarClock, myId)

		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.Request
		msg.SeqNum = append(msg.SeqNum, getValueClock(&scalarClock)[0])
		msg.Text = text
		msg.SendID = myId

		send_to_peer(msg)
	}
}

func send_to_peer(msg utility.Message) {

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utility.Info)
		//open connection whit peer
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
