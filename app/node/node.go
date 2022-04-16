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

	RunTest = true // true if you want to execute test
	peers = list.New()
	utility.Registration(peers, 2345)

	for e := peers.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Info)
		log.Printf("Address: %s:%s", item.Address, item.Port)
	}
	//get myId
	setMyID()
	node := new(Node)
	listNode = append(listNode, *node)
	listNode[0].id = myId
	listNode[0].state = 1
	//fmt.Printf("id in list node %d \n", listNode[0].id)
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
}

func sendmsgWhitdelay(msg string, i int) {

	if !(i == 0) {
		DelayS(GetRandInt(i))
		fmt.Printf("attesa: %d \n", GetRandInt(i))
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

func DelayS(exactTime int) {
	time.Sleep(time.Duration(exactTime) * time.Second)
}

func GetRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(max)
	return n
}
