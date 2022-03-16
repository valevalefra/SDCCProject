package utility

import (
	"errors"
	"log"
	"os"
)

type Utility int

// save registration info to reg_node procedure
func (utils *Utility) Save_registration(arg *Info, res *Result_file) error {
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
	err = checkfile(res)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
