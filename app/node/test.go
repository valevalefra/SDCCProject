package main

import (
	"SDCCProject/app/utility"
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	results = make(map[int]bool)
)

func startTests() {

	//Run tests
	executeTest(1, testLamport)
	//executeTest(2, ricartAgrawala)

}

func executeTest(id int, test func(testId int)) {
	log.Printf("Starting test number %d\n", id)
	test(id)
}

/*
	Testing number of messagges shared between process for Lamport
*/
func testLamport(testId int) {

	msgs := [1]string{"LA-> txt: 'ciao sono il nodo " + strconv.Itoa(myId) + "'"}

	algorithmChoosen = 0
	for _, s := range msgs {
		sendmsgWhitdelay(s, 10)
	}

	time.Sleep(time.Duration(20) * time.Second)
	if utility.MAXPEERS*3 == listNode[0].numberOfMessage {
		log.Printf("Test number %d PASS\n", testId)
	}

}

/*
	Testing number of messagges shared between process for Ricart-Agrawal
*/
func ricartAgrawala(testId int) {

	algorithmChoosen = 1
	msgs := [1]string{"RA-> txt: 'ciao sono il nodo " + strconv.Itoa(myId) + "'"}

	for _, s := range msgs {
		sendmsgWhitdelay(s, 10)
	}

	time.Sleep(time.Duration(40) * time.Second)
	fmt.Printf("list node: %d \n", listNode[0].numberOfMessage)
	if utility.MAXPEERS*2-1 == listNode[0].numberOfMessage {
		log.Printf("Test number %d PASS\n", testId)
	}

}
