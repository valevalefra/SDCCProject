package main

import (
	"SDCCProject/app/utils"
	"container/list"
)

var (
	peers *list.List
)

func main() {

	peers = list.New()
	utils.Registration(peers, 2345)

}
