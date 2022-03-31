package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	peers *list.List
	myId  int
	allId []int
	delay int
)

var listNode []Node

type Node struct {
	id              int
	state           utility.NodeState
	numberOfMessage int
}

var RunTest bool

func main() {

	RunTest = true
	peers = list.New()
	utility.Registration(peers, 2345)

	//const numMsg = 2
	//msgs := []string{"1"}

	for e := peers.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Info)
		log.Printf("Address: %s:%s", item.Address, item.Port)
	}
	//get myId
	setMyID()
	//start clock
	node := Node{
		id:    myId,
		state: 1,
	}

	listNode = append(listNode, node)
	listNode[0].id = 0
	fmt.Println(listNode[0].id)
	startClocks()

	//open listen channel for messages
	//service on port 2345
	//go non bloccante, può continuare a fare altro
	go channel_for_message()

	if RunTest {
		startTests()
		os.Exit(2) //test complete
	}

	menu()

	//	for _, s := range msgs {
	//		fmt.Println("for _, s := range msgs ")
	//		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 100)
	//	}

}

func sendMsg_whitDelay(msg string, i int) {

	if !(delay == 0) {
		Delay_sec(GetRandInt(i))
	}
	err := sendMessages(msg)
	if err != nil {
		return
	}

}

func setMyID() {

	for e := peers.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Info)
		if item.Address == utility.GetLocalIP() {
			myId = item.ID
			allId = append(allId, item.ID)
		} else {
			allId = append(allId, item.ID)
		}
	}
}

func Delay_sec(exactTime int) {
	time.Sleep(time.Duration(exactTime) * time.Second)
}

func GetRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(max)
	return n
}
