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

	const numMsg = 10
	//msgs := [numMsg]string{"ciao", "hello"}

	for e := peers.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Info)
		log.Printf("Address: %s:%s", item.Address, item.Port)
	}
	//get myId
	setMyID()
	//start clock
	clock := utility.ScalarClock{}
	clock.Start()

	//open listen channel for messages
	//service on port 2345
	channel_for_message()

	/*for _, s := range msgs {
		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 2)
	}*/

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
