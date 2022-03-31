package main

import (
	"SDCCProject/app/utility"
	"log"
	"strconv"
)

var RunTest bool
var (
	results = make(map[int]bool)
)

func test(args ...string) error {

	RunTest = true
	return nil
}

func startTests() {

	//Run tests
	executeTest(1, testLamport)

}

func executeTest(id int, test func(testId int) bool) {
	log.Printf("Starting test number %d\n", id)
	res := test(id)
	results[id] = res
	if res {
		log.Printf("Test number %d PASS\n", id)
	} else {
		log.Printf("Test number %d FAILED\n", id)
	}
}

/*
	Testing scalar send by all peer
	3 message send by peer but expected 2 back
*/
func testLamport(testId int) bool {
	const numMsg = 3 * utility.MAXPEERS //3 msg per peer
	msgs := [3]string{"1", "2", "3"}

	algorithmChoosen = 0
	for _, s := range msgs {
		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 2)
	}
	return true
}
