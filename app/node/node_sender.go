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

		if algorithmChoosen == 1 {
			fmt.Printf(" %d  %d %d ", listNode[0].id, listNode[1].id, listNode[2].id)
			if listNode[1].id == myId {
				listNode[myId].state = 2
				fmt.Printf("sono il processo con id %d e ho cambiato il mio stato %d", myId, listNode[myId])
			}
		}

		//send_to_peer(msg, -1)

	}

}

func send_to_peer(msg utility.Message, senderId int) {

	if senderId == -1 {
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
	//send to other peer excluded me
	if senderId == -2 {
		for e := peers.Front(); e != nil; e = e.Next() {
			log.Printf("sto per mandare il mess di release, sono il nodo con id %d \n", msg.SendID)
			if e.Value.(utility.Info).ID != msg.SendID {
				dest := e.Value.(utility.Info)
				log.Printf("sto per mandare il mess di release al nodo con id %d sono il nodo con id %d \n", e.Value.(utility.Info).ID, msg.SendID)
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
	}

	for e := peers.Front(); e != nil; e = e.Next() {
		if e.Value.(utility.Info).ID == senderId {
			dest := e.Value.(utility.Info)
			//Each peer open connection whit peer with sendId
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
}

func send_reply(id int, text string) {
	//prepare msg to send to param id
	var msg utility.Message
	msg.Type = 2
	msg.Text = text
	msg.SendID = myId

	send_to_peer(msg, id)

}

func send_release(msgToDelete utility.Message) {
	//prepare msg to send to other peer
	var msg utility.Message
	msg.Type = 3
	msg.Text = msgToDelete.Text
	msg.SendID = msgToDelete.SendID
	msg.Clock = msgToDelete.Clock

	send_to_peer(msg, -2)

}
