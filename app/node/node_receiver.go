package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

var (
	scalarMsgQueue *list.List
)

func channel_for_message() {
	listener, err := net.Listen("tcp", ":"+"2345")
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()
	scalarMsgQueue = list.New()
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		fmt.Println("CHANNEL FOR MESG")
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {

	defer connection.Close()
	msg := new(utility.Message)
	dec := gob.NewDecoder(connection)
	dec.Decode(msg)

	//update clock
	tmp := msg.SeqNum
	fmt.Println("TMP", tmp)
	updateClock(&scalarClock, tmp)
	incrementClock(&scalarClock, myId)
	fmt.Println("clock", &scalarClock)
	//add in queue and send ack
	e := scalarMsgQueue.PushBack(*msg)
	fmt.Println("PRINT *msg:", *msg)
	fmt.Println("PRINT &msg:", &msg)
	fmt.Println("PRINT Queue:", e.Value)
	//for e := scalarMsgQueue.Front(); e != nil; e = e.Next() {
	//	item := e.Value.(utility.Message)
	//	log.Printf("MESSAGE IN QUEUE: seq num[0] %d: seq num[1] %d :send id %d: ts %d: text %s:tipo %d", item.SeqNum[0], item.SeqNum[1], item.SendID, item.TS, item.Text, item.Type)
	//}
	//e := InsertInOrder(scalarMsgQueue, *msg)
	tmpId := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)
	fmt.Println(tmpId)

	//go scalarMsgDemon(msg, e)
	//go send_scalar_ack(tmpId)

}

func InsertInOrder(l *list.List, msg utility.Message) *list.Element {
	//scan list element for the right position
	tmp := msg.SeqNum[0]
	//fmt.Println("MSG whit seq: "+ strconv.FormatUint(tmp,10))
	for e := l.Front(); e != nil; e = e.Next() {
		item := utility.Message(e.Value.(utility.Message))
		//fmt.Println("ITEM whit seq: "+ strconv.FormatUint(item.SeqNum,10))
		if tmp < item.SeqNum[0] {
			//found the next item
			//fmt.Println("IF CONDITION OK")
			return l.InsertBefore(msg, e)
		}
	}
	//fmt.Println("PUSHBACK")
	return l.PushBack(msg)
}
