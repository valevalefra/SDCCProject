package utility

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Result_file struct {
	PeerNum int
	Peers   [MAXPEERS]string
}

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

type NodeState int

const (
	cs  NodeState = 0
	ncs           = 1
	req           = 2
)

// Struct to send information about peer
type Info struct {
	ID      int
	Address string
	Port    string
	State   NodeState //used for Ricart-Agrawala's algorithm
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func setInfo(info *Info, port int) error {
	info.Address = GetLocalIP()
	if info.Address == "" {
		return errors.New("Impossible to find local ip")
	}

	info.Port = strconv.Itoa(port)
	info.State = 1 // initially all node are in ncs (not critical section)
	return nil
}

func (n *Info) ChangeState(i int) error {
	n.State = NodeState(i)
	return nil
}

func ParseLine(s string, sep string) (string, string) {
	res := strings.Split(s, sep)
	return res[0], res[1]
}

func checkfile(res *Result_file) error {
	res.PeerNum = MAXPEERS
	file, err := os.Open("/tmp/clients.txt")
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		res.Peers[i] = line
		i++
	}
	return nil
}

/*
	Registration function for peer
*/
func Registration(peers *list.List, port int) {

	var info Info
	var res Result_file

	addr := Server_addr + ":" + strconv.Itoa(Server_port)
	// each peer Try to connect to addr of registry
	server, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer server.Close()

	err = setInfo(&info, port)
	if err != nil {
		log.Fatal("Error on setInfo: ", err)
	}

	//call procedure
	log.Printf("Call to registration node")
	err = server.Call("Utility.SaveRegistration", &info, &res)
	if err != nil {
		log.Fatal("Error save_registration procedure: ", err)
	}

	//check result
	for e := 0; e < MAXPEERS; e++ {
		var item Info
		item.Address, item.Port = ParseLine(res.Peers[e], ":")
		item.ID = e
		item.State = 1
		peers.PushBack(item)

	}

}
