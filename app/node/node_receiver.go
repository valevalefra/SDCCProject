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
		fmt.Printf("Prima [%s]: %d\n", text, ackCounter[text])
		mutex.Lock()
		ackCounter[text] = ackCounter[text] + 1
		mutex.Unlock()
		fmt.Printf("Dopo [%s]: %d\n ", text, ackCounter[text])
	}

}

func handleConnection(connection net.Conn) {

	//defer connection.Close()
	msg := new(utility.Message)
	dec := gob.NewDecoder(connection)
	dec.Decode(msg)

	switch msg.Type {
	case utility.Request:
		if algorithmChoosen == 0 {
			fmt.Printf("algo scelto %d \n", algorithmChoosen)
			//update clock
			tmp := msg.Clock[0]
			fmt.Printf("il nodo con id %d ha ricevuto il messaggio di richiesta %s, che ha valore del clock tmp %d, dal nodo con id %d (i due id dovrebbero essere uguali) \n", myId, msg.Text, tmp, msg.SendID)
			updateClock(&scalarClock, tmp)
			fmt.Printf("il nodo con id %d ha fatto update del clock, il valore del clock ora è %d \n", myId, *&scalarClock)
			incrementClock(&scalarClock)
			fmt.Printf("il nodo con id %d incrementa il valore del clock di una unità  %d \n", myId, *&scalarClock)
			//add in queue and send ack
			//scalarMsgQueue.PushBack(*msg)
			//fmt.Println("PRINT *msg:", *msg) printa contenuto mess
			//fmt.Println("PRINT &msg:", &msg) printa indirizzo di memoria
			//fmt.Println("PRINT Queue:", e.Value)
			e := Reordering(scalarMsgQueue, *msg)
			fmt.Println("PRINT Queue:", e.Value)
			//tmpId := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)
			//fmt.Println(tmpId)
			//for lamport
			fmt.Println("prima di controllo condizione")
			go checkCondition(msg, e)
			fmt.Println("dopo di controllo condizione")
			go send_reply(msg.SendID, msg.Text)
		} else {
			tmp := msg.Clock[0]
			fmt.Printf("il nodo con id %d ha ricevuto il messaggio di richiesta %s, che ha valore del clock tmp %d, dal nodo con id %d (i due id dovrebbero essere uguali) \n", myId, msg.Text, tmp, msg.SendID)
			updateClock(&scalarClock, tmp)
			go check_numberOfReply(msg)
			go reply_and_check(scalarMsgQueue, *msg)
		}

	case utility.Reply:

		fmt.Printf("il nodo con id %d ha ricevuto un mess di reply, per il mess %s, dal nodo con id %d \n", myId, msg.Text, msg.SendID)
		text := msg.Text
		//fmt.Println("ACK FOR: " + text)
		ackChan <- text

	case utility.Release:
		if algorithmChoosen == 0 {
			fmt.Printf("il nodo con id %d ha ricevuto un mess di release, per il mess %s, dal nodo con id %d \n", myId, msg.Text, msg.SendID)
			//case release cancella messaggio dalla coda.
			l := scalarMsgQueue
			if l.Len() != 0 {
				fmt.Printf("lunghezza coda in switch release %d\n", l.Len())
				for e := l.Front(); e != nil; e = e.Next() {
					if e.Value.(utility.Message).SendID == msg.SendID && e.Value.(utility.Message).Clock[0] == msg.Clock[0] {
						fmt.Printf("il nodo con id %d ha ricevuto un mess di release quindi sta elimando dalla propria coda il mess %s \n", myId, e.Value.(utility.Message).Text)
						scalarMsgQueue.Remove(e)
					}
				}
				fmt.Printf("lunghezza coda in switch release DOPO RIMOZIONE %d\n", l.Len())
			}
		}
	}
}

//function for ricart agrawala, check number of replies for request
func check_numberOfReply(msg *utility.Message) {
	if scalarMsgQueue.Len() != 0 {
		for !(count_reply(*msg)) {
			utility.Delay_ms(100)
		}
		listNode[0].state = 0 //TODO: casomai simula tempo più lungo per la sezione critica
		fmt.Println("condizione verificata, puoi accedere alla sezione critica \n")
		enterCS(scalarMsgQueue.Front().Value.(utility.Message))
		//delete msg from my queue
		fmt.Printf("rimuovo dalla coda il mess del processo %d il cui testo era %s \n", msg.SendID, scalarMsgQueue.Front().Value.(utility.Message).Text)
		msgToDelete := scalarMsgQueue.Front().Value.(utility.Message)
		msgID := strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).SendID) + "-" + strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).Clock[0])
		delete(ackCounter, msgID)
		fmt.Printf("rimuovo ackcounter per il mess %s \n", msgID)
		scalarMsgQueue.Remove(scalarMsgQueue.Front())
		listNode[0].state = 1
		send_release_to(msgToDelete, scalarMsgQueue)
	}

}

func count_reply(message utility.Message) bool {
	if scalarMsgQueue.Len() != 0 {
		//get first element on queue
		tmp := scalarMsgQueue.Front().Value.(utility.Message)
		//tmpId := strconv.Itoa(tmp.SendID) + "-" + strconv.Itoa(tmp.Clock[0])
		//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
		//fmt.Println("tmpid:", tmpId, " num ack: ", ackCounter[tmp.Text])
		mutex.Lock()
		//forse andrebbe modificato l'identificativo
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
func reply_and_check(queue *list.List, msg utility.Message) {

	// se è in sc mette mess in coda
	if listNode[0].state == 0 {
		e := Reordering(queue, msg)
		fmt.Printf("sono il processo %d sono in sezione critica quindi metto %s in coda, la mia coda sarà %s, lunghezza coda %d \n", myId, *&msg.Text, e.Value, queue.Len())

	}
	//se è interessato alla sc mette mess in coda
	c := getValueClock(&scalarClock)
	fmt.Printf("sono il processo %d e il mio clock in reply and check è %d \n", myId, c[0])
	if listNode[0].state == 2 && c[0] <= msg.Clock[0] || listNode[0].state == 1 && c[0] == msg.Clock[0] && myId <= msg.SendID {
		e := Reordering(queue, msg)
		fmt.Printf("sono il processo %d sono in req per sc quindi metto %s in coda, la mia coda sarà %s, lunghezza coda %d \n", myId, *&msg.Text, e.Value, queue.Len())
	}
	// se non è interessato alla sc e non è in sc allora manda il reply al processo con id: msg.sendID
	if listNode[0].state == 1 {
		fmt.Printf("sono il processo %d non sono in sc quindi invio %s al processo con id %d \n", myId, *&msg.Text, msg.SendID)
		send_reply(msg.SendID, msg.Text)
	}
}

func checkCondition(msg *utility.Message, e *list.Element) {

	//first condition
	if scalarMsgQueue.Len() != 0 {
		for !(firstCondition(*msg) && !secondCondition(*msg)) {
			utility.Delay_ms(100)
		}
		fmt.Println("prima e seconda condizione verificata, puoi accedere alla sezione critica \n")
		enterCS(scalarMsgQueue.Front().Value.(utility.Message))
		//delete msg from my queue
		fmt.Printf("rimuovo dalla coda il mess del processo %d il cui testo era %s \n", msg.SendID, scalarMsgQueue.Front().Value.(utility.Message).Text)
		//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
		msgToDelete := scalarMsgQueue.Front().Value.(utility.Message)
		msgID := strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).SendID) + "-" + strconv.Itoa(scalarMsgQueue.Front().Value.(utility.Message).Clock[0])
		delete(ackCounter, msgID)
		fmt.Printf("rimuovo ackcounter per il mess %s \n", msgID)
		////////////////////////////////////////////////////////////////////////////////
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
	//path_file := "/docker/node_volume/" + p2.Peer.Ip_address + "_log.txt"
	//f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if _, err := os.Stat("/home/valentina/GolandProjects/SDCCProject/app/data.txt"); errors.Is(err, os.ErrNotExist) {
	fmt.Println("STO SCRIVENDOOOOOOOO")
	//f, err := os.Open("sharedFile")

	if err != nil {
		log.Fatal(err)
	}

	//defer f.Close()

	_, err2 := f.WriteString(message.Text + " " + strconv.Itoa(message.Clock[0]) + " " + strconv.Itoa(myId) + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	//}
	fmt.Println("done")

}

func firstCondition(msg utility.Message) bool {

	if scalarMsgQueue.Len() != 0 {
		//get first element on queue
		tmp := scalarMsgQueue.Front().Value.(utility.Message)
		//tmpId := strconv.Itoa(tmp.SendID) + "-" + strconv.Itoa(tmp.Clock[0])
		//msgID := strconv.Itoa(msg.SendID) + "-" + strconv.Itoa(msg.Clock[0])
		//fmt.Println("tmpid:", tmpId, " num ack: ", ackCounter[tmp.Text])
		mutex.Lock()
		//forse andrebbe modificato l'identificativo
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
