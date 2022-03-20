package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"log"
	"math/rand"
	"time"
)

var (
	peers *list.List
	myId  int
	allId []int
	delay int
)

func main() {

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
	startClocks()

	//open listen channel for messages
	//service on port 2345
	//go non bloccante, può continuare a fare altro
	go channel_for_message()

	menu()

	//	for _, s := range msgs {
	//		fmt.Println("for _, s := range msgs ")
	//		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 100)
	//	}
	for e := scalarMsgQueue.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Message)
		log.Printf("MESSAGE IN QUEUE:send id %d:: text %s:tipo %d", item.SendID, item.Text, item.Type)
	}

}

func sendMsg_whitDelay(msg string, i int) {

	if !(delay == 0) {
		Delay_sec(GetRandInt(delay))
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
