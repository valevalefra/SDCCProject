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

const (
	Server_port int    = 4321
	Server_addr string = "10.10.1.50"
)

type Result_file struct {
	PeerNum int
	Peers   [3]string
}

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

// Struct to send information about peer
type Info struct {
	ID      int
	Address string
	Port    string
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
	return nil
}

func ParseLine(s string, sep string) (string, string) {
	res := strings.Split(s, sep)
	return res[0], res[1]
}

func checkfile(res *Result_file) error {
	res.PeerNum = 3
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
	// Try to connect to addr
	server, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer server.Close()

	//set info to send
	err = setInfo(&info, port)
	if err != nil {
		log.Fatal("Error on setInfo: ", err)
	}

	//call procedure
	log.Printf("Call to registration node")
	err = server.Call("Utility.Save_registration", &info, &res)
	if err != nil {
		log.Fatal("Error save_registration procedure: ", err)
	}

	//check result
	for e := 0; e < 3; e++ {
		var item Info
		item.Address, item.Port = ParseLine(res.Peers[e], ":")
		item.ID = e
		peers.PushBack(item)

	}

}
