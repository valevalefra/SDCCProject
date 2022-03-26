package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

const (
	buffSize int = 1000
)

var (
	scalarMsgQueue *list.List
	ackChan        = make(chan string, buffSize)
	ackCounter     map[string]int //key-value : msg.txt-num ack
	mutex          sync.Mutex
)

func channel_for_message() {
	listener, err := net.Listen("tcp", ":"+"2345")
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	go check_reply()

	scalarMsgQueue = list.New()
	ackCounter = make(map[string]int)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go handleConnection(connection)
	}
}

func check_reply() {

	for text := range ackChan {
		fmt.Printf("Prima [%s]: %d\n", text, ackCounter[text])
		mutex.Lock()
		ackCounter[text] = ackCounter[text] + 1
		mutex.Unlock()
		fmt.Printf("Dopo [%s]: %d\n ", text, ackCounter[text])
	}

}

func handleConnection(connection net.Conn) {

	defer connection.Close()
	msg := new(utility.Message)
	dec := gob.NewDecoder(connection)
	dec.Decode(msg)
	//fmt.Printf("il nodo con id %d e valore del clock %d sta ricevendo %s \n", myId, *&scalarClock, msg.Text)

	switch msg.Type {
	case utility.Request:
		//update clock
		tmp := msg.Clock[0]
		fmt.Printf("il nodo con id %d ha ricevuto il messaggio di richiesta %s, che ha valore del clock tmp %d, dal nodo con id %d (i due id dovrebbero essere uguali) \n", myId, msg.Text, tmp, msg.SendID)
		updateClock(&scalarClock, tmp)
		fmt.Printf("il nodo con id %d ha fatto update del clock, il valore del clock ora è %d \n", myId, *&scalarClock)
		incrementClock(&scalarClock)
		fmt.Printf("il nodo con id %d incrementa il valore del clock di una unità  %d \n", myId, *&scalarClock)
		//add in queue and send ack
		//e := scalarMsgQueue.PushBack(*msg)
		//fmt.Println("PRINT *msg:", *msg) printa contenuto mess
		//fmt.Println("PRINT &msg:", &msg) printa indirizzo di memoria
		//fmt.Println("PRINT Queue:", e.Value)
		e := Reordering(scalarMsgQueue, *msg)
		fmt.Println("PRINT Queue:", e.Value)
		//tmpId := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)
		//fmt.Println(tmpId)

		go checkCondition(msg, e)
		go send_reply(msg.SendID, msg.Text)

	case utility.Reply:

		fmt.Printf("il nodo con id %d ha ricevuto un mess di reply, per il mess %s, dal nodo con id %d \n", myId, msg.Text, msg.SendID)
		text := msg.Text
		//fmt.Println("ACK FOR: " + text)
		ackChan <- text

	case utility.Release:
		fmt.Printf("il nodo con id %d ha ricevuto un mess di release, per il mess %s, dal nodo con id %d \n", myId, msg.Text, msg.SendID)
		//case release cancella messaggio dalla coda.

	}
}

func checkCondition(msg *utility.Message, e *list.Element) {

	//first condition
	if firstCondition(*msg) {
		fmt.Println("prima condizione verificata \n")
	}
	if !secondCondition(*msg) {
		fmt.Println("seconda condizione verificata \n")
	}
	for !(firstCondition(*msg) && !secondCondition(*msg)) {
		utility.Delay_ms(100)
	}
	fmt.Println("prima e seconda condizione verificata, puoi accedere alla sezione critica \n")
	enterCS(scalarMsgQueue.Front().Value.(utility.Message))
	//delete msg from my queue
	fmt.Printf("rimuovo dalla coda il mess del processo %d il cui testo era %s \n", msg.SendID, scalarMsgQueue.Front().Value.(utility.Message).Text)
	fmt.Printf("messaggio element %s \n", e.Value)
	// scalarMsgQueue.Remove(e)
	//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
	msgToDelete := scalarMsgQueue.Front().Value.(utility.Message)
	msgID := strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).SendID) + "-" + strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).Clock[0])
	delete(ackCounter, msgID)
	fmt.Printf("rimuovo ackcounter per il mess %s \n", msgID)
	scalarMsgQueue.Remove(scalarMsgQueue.Front())
	send_release(msgToDelete)
}

func enterCS(message utility.Message) {

	fmt.Println("scrivi su file " + message.Text)

}

func firstCondition(msg utility.Message) bool {

	//get first element on queue
	tmp := scalarMsgQueue.Front().Value.(utility.Message)
	//tmpId := strconv.Itoa(tmp.SendID) + "-" + strconv.Itoa(tmp.Clock[0])
	//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
	//fmt.Println("tmpid:", tmpId, " num ack: ", ackCounter[tmp.Text])
	mutex.Lock()
	ack := ackCounter[tmp.Text]
	mutex.Unlock()
	if ack == utility.MAXPEERS {

		return true
	} else {
		return false
	}
}

func secondCondition(msg utility.Message) bool {

	msgHead := scalarMsgQueue.Front()
	msgFirst := msgHead.Value.(utility.Message)
	if msgFirst.Clock[0] == msg.Clock[0] && msgFirst.SendID == msg.SendID {
		return false
	}
	for i := 0; i < len(allId); i++ {
		check := false
		for e := msgHead.Next(); e != nil; e = e.Next() {
			item := e.Value.(utility.Message)
			fmt.Printf("il messaggio item é %s \n", item.Text)
			fmt.Printf("item.sendID == %d and allID[i]== %d item.clock= %d e msg.clock =%d e la i vale %d \n", item.SendID, allId[i], item.Clock[0], msg.Clock[0], i)
			if item.SendID == allId[i] && item.Clock[0] > msg.Clock[0] {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}

	return true
}

func Reordering(l *list.List, msg utility.Message) *list.Element {
	//scan list element for the right position
	tmp := msg.Clock
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Message).Clock
		if tmp[0] < item[0] {
			//found the next item
			fmt.Println("IF CONDITION OK")
			return l.InsertBefore(msg, e)
		}
		if tmp[0] == item[0] {
			//found the next item
			tmp := msg.SendID
			idFirst := e.Value.(utility.Message).SendID
			fmt.Println("IF CONDITION SONO UGUALI I VALORI DEI CLOCK")
			if tmp < idFirst {
				return l.InsertBefore(msg, e)
			}

		}
	}
	return l.PushBack(msg)
}
