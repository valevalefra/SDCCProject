package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func sendMessages(args ...string) error {
	var function func(msgs []string)
	function = sendTo
	go function(args)
	return nil

}

func sendTo(msgs []string) {

	for _, text := range msgs {
		//increment local clock
		incrementClock(&scalarClock)
		fmt.Printf("il nodo con id %d e valore del clock %d sta inviando %s \n", myId, scalarClock.GetValue(), text)

		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.Request
		msg.Clock = getValueClock(&scalarClock)
		msg.Text = text
		msg.SendID = myId

		send_to_peer(msg, -1)

	}

}

func send_to_peer(msg utility.Message, senderId int) {
	//Send to all peer
	if senderId == -1 {
		for e := peers.Front(); e != nil; e = e.Next() {
			listNode[0].numberOfMessage = listNode[0].numberOfMessage + 1
			dest := e.Value.(utility.Info)
			//open connection whit other peer
			peer_conn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peer_conn)

			if err != nil {
				log.Println("Send response error on Dial")
				for err != nil {
					time.Sleep(time.Duration(rand.Int()))
					conn, err = net.Dial("tcp", peer_conn)
				}
			}
			//Ricart Agrawala
			if algorithmChoosen == 1 {
				if listNode[0].id == myId {
					listNode[0].state = 2 //set state of peer to requesting (cs)
					fmt.Printf("sono il processo con id %d e ho cambiato il mio stato in %d (request cs) \n", myId, listNode[0].state)
				}
			}
			enc := gob.NewEncoder(conn)
			error := enc.Encode(msg)
			if error != nil {
				log.Println("Error in encoder")
			}
			//defer conn.Close()
		}
	}
	//send to other peer excluded me
	if senderId == -2 {
		for e := peers.Front(); e != nil; e = e.Next() {
			listNode[0].numberOfMessage = listNode[0].numberOfMessage + 1
			if e.Value.(utility.Info).ID != msg.SendID {
				dest := e.Value.(utility.Info)
				log.Printf("sto per mandare il mess di release al nodo con id %d, sono il nodo con id %d \n", e.Value.(utility.Info).ID, msg.SendID)
				//open connection whit other peer
				peer_conn := dest.Address + ":" + dest.Port
				conn, err := net.Dial("tcp", peer_conn)
				if err != nil {
					log.Println("Send response error on Dial")
					for err != nil {
						time.Sleep(time.Duration(rand.Int()))
						conn, err = net.Dial("tcp", peer_conn)
					}
				}
				enc := gob.NewEncoder(conn)
				enc.Encode(msg)
				//defer conn.Close()
			}
		}
	}
}

func send_reply(id int, text string) {
	//prepare msg to send to param id
	var msg utility.Message
	msg.Type = 2
	msg.Text = text
	msg.SendID = myId
	listNode[0].numberOfMessage = listNode[0].numberOfMessage + 1
	//send to specific peer
	for e := peers.Front(); e != nil; e = e.Next() {
		if e.Value.(utility.Info).ID == id {
			dest := e.Value.(utility.Info)
			//Each peer open connection whit peer with sendId
			peer_conn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peer_conn)
			for err != nil {
				time.Sleep(time.Duration(rand.Int()))
				conn, err = net.Dial("tcp", peer_conn)
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)
			//defer conn.Close()
		}

	}

}

func send_release(msgToDelete utility.Message) {
	//prepare msg to send to other peer
	var msg utility.Message
	msg.Type = 3
	msg.Text = msgToDelete.Text
	msg.SendID = msgToDelete.SendID
	msg.Clock = msgToDelete.Clock
	fmt.Println("dentro send release")

	send_to_peer(msg, -2)

}

//function for ricart agrawala. peer send release only to process in own queue
func send_release_to(toDelete utility.Message, l *list.List) {

	//prepare msg to send to other peer
	var msg utility.Message
	//msg.Type = 2 //reply
	//msg.Text = toDelete.Text
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Message).SendID
		mess := e.Value.(utility.Message).Text
		msg.Type = 2 //release
		msg.Text = mess
		msg.SendID = item
		scalarMsgQueue.Remove(e)
		send_reply(item, mess)

	}

}
