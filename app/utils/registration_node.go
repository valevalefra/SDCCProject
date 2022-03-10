package utils

import (
	"bufio"
	"container/list"
	"errors"
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

// save registration info to reg_node procedure
func Save_registration(arg *Info, res *Result_file) error {
	log.Printf("The registration is for node whith ip address:port : %s:%s\n", arg.Address, arg.Port)
	f, err := os.OpenFile("/tmp/clients.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		return errors.New("Impossible to open file")
	}
	/*
		see https://www.joeshaw.org/dont-defer-close-on-writable-files/ for file defer on close
	*/
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	//save new address on file
	_, err = f.WriteString(arg.Address + ":" + arg.Port)
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		return err
	}

	log.Printf("Saved")

	Connection <- true
	Wg.Add(1)
	log.Printf("Waiting other connection")
	Wg.Wait()

	//send back file
	err = prepare_response(res)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func prepare_response(res *Result_file) error {
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
		res.Peers[i] = line
		i++
	}
	if err := scanner.Err(); err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	err = file.Sync()
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	return nil
}

/*
	Registration function for peer
*/
func Registration(peers *list.List, port int) {

	var info Info
	var res Result_file

	addr := "10.10.1.50" + ":" + "4321"
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
	err = server.Call("Save_registration", &info, &res)
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
