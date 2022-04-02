package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
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

type Algorithm int

const (
	lmp Algorithm = 0
	ra            = 1
)

var algorithmChoosen Algorithm

func channel_for_message() {
	listener, err := net.Listen("tcp", ":"+"2345")
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	//defer listener.Close()

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
		//fmt.Printf("Prima [%s]: %d\n", text, ackCounter[text])
		mutex.Lock()
		ackCounter[text] = ackCounter[text] + 1
		mutex.Unlock()
		//fmt.Printf("Dopo [%s]: %d\n ", text, ackCounter[text])
	}

}

func handleConnection(connection net.Conn) {

	//defer connection.Close()
	msg := new(utility.Message)
	dec := gob.NewDecoder(connection)
	err := dec.Decode(msg)
	if err != nil {
		return
	}

	switch msg.Type {
	case utility.Request:
		if algorithmChoosen == 0 {
			//update clock
			tmp := msg.Clock[0]
			fmt.Printf("il nodo con id %d ha ricevuto il messaggio di richiesta %s, che ha valore del clock: %d, dal nodo con id %d \n", myId, msg.Text, tmp, msg.SendID)
			updateClock(&scalarClock, tmp)
			//fmt.Printf("il nodo con id %d ha fatto update del clock, il valore del clock ora è %d \n", myId, *&scalarClock)
			incrementClock(&scalarClock)
			//fmt.Printf("il nodo con id %d incrementa il valore del clock di una unità  %d \n", myId, *&scalarClock)
			//add in queue and send ack
			Reordering(scalarMsgQueue, *msg)
			//for lamport
			go checkCondition(msg)
			go send_reply(msg.SendID, msg.Text)
		} else {
			tmp := msg.Clock[0]
			fmt.Printf("il nodo con id %d ha ricevuto il messaggio di richiesta %s, che ha valore del clock %d, dal nodo con id %d  \n", myId, msg.Text, tmp, msg.SendID)
			updateClock(&scalarClock, tmp)
			go checkNumberofreply()
			go replyAndCheck(scalarMsgQueue, *msg)
		}

	case utility.Reply:

		fmt.Printf("il nodo con id %d ha ricevuto un mess di reply, per il mess %s, dal nodo con id %d che ha valore del clock pari a %d \n", myId, msg.Text, msg.SendID, msg.Clock[0])
		text := msg.Text
		ackChan <- text

	case utility.Release:
		if algorithmChoosen == 0 {
			fmt.Printf("il nodo con id %d ha ricevuto un mess di release, per il mess %s, dal nodo con id %d che ha valore del clock pari a %d \n", myId, msg.Text, msg.SendID, msg.Clock[0])
			//delete msg from queue
			l := scalarMsgQueue
			if l.Len() != 0 {
				for e := l.Front(); e != nil; e = e.Next() {
					if e.Value.(utility.Message).SendID == msg.SendID && e.Value.(utility.Message).Clock[0] == msg.Clock[0] {
						fmt.Printf("il nodo con id %d ha ricevuto un mess di release quindi sta elimando dalla propria coda il mess %s \n", myId, e.Value.(utility.Message).Text)
						scalarMsgQueue.Remove(e)
					}
				}
			}
		}
	}
}

///////////////////////////// RICART AGRAWALA ///////////////////////////////////////////////

//function for ricart agrawala, check number of replies for request
func checkNumberofreply() {
	if scalarMsgQueue.Len() != 0 {
		for !(countReply()) {
			utility.Delay_ms(100)
		}
		listNode[0].state = 0
		fmt.Printf("condizione verificata, puoi accedere alla sezione critica il mio stato è %d \n", listNode[0].state)
		enterCS(scalarMsgQueue.Front().Value.(utility.Message))
		//delete msg from my queue
		fmt.Printf("rimuovo dalla coda il mess del processo %d il cui testo era %s \n", scalarMsgQueue.Front().Value.(utility.Message).SendID, scalarMsgQueue.Front().Value.(utility.Message).Text)
		msgToDelete := scalarMsgQueue.Front().Value.(utility.Message)
		msgID := strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).SendID) + "-" + strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).Clock[0])
		delete(ackCounter, msgID)
		fmt.Printf("rimuovo ackcounter per il mess %s \n", msgID)
		l := scalarMsgQueue
		if l.Len() != 0 {
			for e := l.Front(); e != nil; e = e.Next() {
				if e.Value.(utility.Message).SendID == msgToDelete.SendID && e.Value.(utility.Message).Clock[0] == msgToDelete.Clock[0] {
					scalarMsgQueue.Remove(e)
				}
			}
		}
		fmt.Printf("lunghezza coda DOPO RIMOZIONE: %d \n", l.Len())
		listNode[0].state = 1
		fmt.Printf("uscito dalla sc,il mio stato è %d, mando mess di release ai nodi nella mia coda \n", listNode[0].state)
		send_release_to(msgToDelete, scalarMsgQueue)
	}

}

func countReply() bool {
	if scalarMsgQueue.Len() != 0 {
		//get first element on queue
		tmp := scalarMsgQueue.Front().Value.(utility.Message)
		mutex.Lock()
		ack := ackCounter[tmp.Text]
		mutex.Unlock()
		if ack == (utility.MAXPEERS - 1) {

			return true
		} else {
			return false
		}
	}
	return false

}

//function for ricart agrawala menage request
func replyAndCheck(queue *list.List, msg utility.Message) {

	// if node's state is "cs" put message in queue
	if listNode[0].state == 0 {
		Reordering(queue, msg)
		fmt.Printf("sono il processo %d sono in sezione critica quindi metto il messaggio %s in coda, lunghezza coda %d \n", myId, *&msg.Text, queue.Len())

	}
	c := getValueClock(&scalarClock)
	//fmt.Printf("sono il processo %d e il mio clock in reply and check è %d \n", myId, c[0])
	// if node's state is "request cs" and its clock is lower than other node or its node's id is lower than other node, then insert msg in queue
	if listNode[0].state == 2 && c[0] < msg.Clock[0] || listNode[0].state == 2 && c[0] == msg.Clock[0] && myId <= msg.SendID {
		Reordering(queue, msg)
		fmt.Printf("sono il processo %d sono nello stato di request per cs quindi metto il messaggio %s in coda, lunghezza coda %d \n", myId, *&msg.Text, queue.Len())
	}
	// if node's state is "request cs" and its clock is higher than other node or its node's id is higher than other node, then send reply to the node that has the right to access
	if listNode[0].state == 2 && c[0] > msg.Clock[0] || listNode[0].state == 2 && c[0] == msg.Clock[0] && myId > msg.SendID {
		//e := Reordering(queue, msg)
		fmt.Printf("sono il processo %d sono nello stato di request per cs ma non ho diritto alla cs quindi mando reply a %d, lunghezza coda %d \n", myId, msg.SendID, queue.Len())
		send_reply(msg.SendID, msg.Text)
	}

	// if node's state is "ncs" send reply to node wih id: msg.sendID
	if listNode[0].state == 1 {
		fmt.Printf("sono il processo %d non sono in sc e non sono interessato ad accedere quindi invio %s al processo con id %d \n", myId, *&msg.Text, msg.SendID)
		send_reply(msg.SendID, msg.Text)
	}
}

///////////////////////////////////// LAMPORT ///////////////////////////////////////////////

func checkCondition(msg *utility.Message) {

	//first condition
	if scalarMsgQueue.Len() != 0 {
		for !(firstCondition() && !secondCondition(*msg)) {
			utility.Delay_ms(100)
		}
		fmt.Println("prima e seconda condizione verificata, puoi accedere alla sezione critica")
		enterCS(scalarMsgQueue.Front().Value.(utility.Message))
		//delete msg from my queue
		fmt.Printf("rimuovo dalla coda il mess del processo %d il cui testo era %s \n", msg.SendID, scalarMsgQueue.Front().Value.(utility.Message).Text)
		//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
		msgToDelete := scalarMsgQueue.Front().Value.(utility.Message)
		msgID := strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).SendID) + "-" + strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).Clock[0])
		delete(ackCounter, msgID)
		fmt.Printf("rimuovo ackcounter per il mess %s \n", msgID)
		l := scalarMsgQueue
		if l.Len() != 0 {
			fmt.Printf("lunghezza coda: %d \n", l.Len())
			for e := l.Front(); e != nil; e = e.Next() {
				if e.Value.(utility.Message).SendID == msgToDelete.SendID && e.Value.(utility.Message).Clock[0] == msgToDelete.Clock[0] {
					fmt.Printf("il nodo con id %d ha ricevuto un mess di release quindi sta elimando dalla propria coda il mess %s \n", myId, e.Value.(utility.Message).Text)
					scalarMsgQueue.Remove(e)
				}
			}
			fmt.Printf("lunghezza coda DOPO RIMOZIONE: %d \n", l.Len())
		}
		send_release(msgToDelete)

	}
}

func enterCS(message utility.Message) {

	fmt.Println("scrivi su file " + message.Text)
	path_file := "/docker/node_volume/" + "log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Println("STO SCRIVENDO IL MESSAGGIO")

	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	_, err2 := f.WriteString(message.Text + " " + strconv.Itoa(message.Clock[0]) + " " + strconv.Itoa(myId) + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("done")

}

func firstCondition() bool {

	if scalarMsgQueue.Len() != 0 {
		//get first element on queue
		tmp := scalarMsgQueue.Front().Value.(utility.Message)
		mutex.Lock()
		ack := ackCounter[tmp.Text]
		mutex.Unlock()
		if ack == utility.MAXPEERS {

			return true
		} else {
			return false
		}
	}
	return false
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
			//fmt.Printf("il messaggio item é %s \n", item.Text)
			//fmt.Printf("item.sendID == %d and allID[i]== %d item.clock= %d e msg.clock =%d e la i vale %d \n", item.SendID, allId[i], item.Clock[0], msg.Clock[0], i)
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
			return l.InsertBefore(msg, e)
		}
		if tmp[0] == item[0] {
			//found the next item
			tmp := msg.SendID
			idFirst := e.Value.(utility.Message).SendID
			fmt.Println("valori dei clock uguali")
			if tmp < idFirst {
				return l.InsertBefore(msg, e)
			}

		}
	}
	return l.PushBack(msg)
}
