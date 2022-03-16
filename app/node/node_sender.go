package main

import (
	"SDCCProject/app/utility"
	"encoding/gob"
	"log"
	"net"
	"time"
)

func sendMessages(msgs []string) error {
	for _, text := range msgs {
		//increment local clock
		incrementClock(&scalarClock, myId)

		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.Request
		msg.SeqNum = append(msg.SeqNum, getValueClock(&scalarClock)[0])
		msg.Date = time.Now().Format("2006/01/02 15:04:05")
		msg.Text = text
		msg.SendID = myId

		send_to_peer(msg)
	}
	return nil
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
