package main

import (
	"SDCCProject/app/prova"
	"container/list"
)

var (
	peers *list.List
)

func main() {

	peers = list.New()
	prova.Proof()
	//utils.Registration(peers, 2345)

}
