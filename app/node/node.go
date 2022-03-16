package main

import (
	"SDCCProject/app/utility"
	"container/list"
	"fmt"
	"log"
)

var (
	peers *list.List
	myId  int
	allId []int
)

func main() {

	peers = list.New()
	utility.Registration(peers, 2345)

	for e := peers.Front(); e != nil; e = e.Next() {
		item := e.Value.(utility.Info)
		log.Printf("Address: %s:%s", item.Address, item.Port)
	}
	//get myId
	setMyID()
	//start clock
	//startClocks()

	//open listen channel for messages
	//service on port 2345
	//go message_handler()

}

func setMyID() {

	for e := peers.Front(); e != nil; e = e.Next() {
		fmt.Errorf("ciaooooooo", e)
		item := e.Value.(utility.Info)
		if item.Address == utility.GetLocalIP() {
			myId = item.ID
			allId = append(allId, item.ID)
		} else {
			allId = append(allId, item.ID)
		}
	}
}
