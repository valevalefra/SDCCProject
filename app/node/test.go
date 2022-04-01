package main

import (
	"SDCCProject/app/utility"
	"log"
	"time"
)

var (
	results = make(map[int]bool)
)

func startTests() {

	//Run tests
	executeTest(1, testLamport)

}

func executeTest(id int, test func(testId int)) {
	log.Printf("Starting test number %d\n", id)
	test(id)
	//res := test(id)
	/*results[id] = res
	if res {
		log.Printf("Test number %d PASS\n", id)
	} else {
		log.Printf("Test number %d FAILED\n", id)
	}*/
}

/*
	Testing scalar send by all peer
	3 message send by peer but expected 2 back
*/
func testLamport(testId int) {
	const numMsg = 3 * utility.MAXPEERS //3 msg per peer
	msgs := [1]string{"abc"}

	algorithmChoosen = 0
	for _, s := range msgs {
		sendMsg_whitDelay(s, 10)
	}

	time.Sleep(time.Duration(20) * time.Second)

}
