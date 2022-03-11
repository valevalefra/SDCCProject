package main

import (
	"SDCCProject/app/utility"
	"container/list"
)

var (
	peers *list.List
)

func main() {

	peers = list.New()
	utility.Registration(peers, 2345)

}
