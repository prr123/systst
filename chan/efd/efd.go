//https://github.com/gxed/eventfd/tree/80a92cca79a8041496ccc9dd773fcb52a57ec6f9
package main

import (
	"fmt"
	"github.com/sahne/eventfd"
	"syscall"
)

func main() {
	var ba []byte

	var val uint64

	efd, err := eventfd.New()
	if err != nil {
		fmt.Println("Could not create EventFD: %v", err)
		return
	}

	eFD := efd.Fd()
	fmt.Println("fd: ", eFD)

	ba = []byte{2,0,0,0,0,0,0,0}
	pid := syscall.Getpid()


	fmt.Println("pid: ", pid," ba: ", ba, len(ba))

	n, err := syscall.Write(eFD, ba)
	if err != nil {
		fmt.Println("error efd syscall write: ", err)
		return
	}
	fmt.Println("write:  ", n)


	val = 0
/*
	err = efd.WriteEvents(val)
	if err !=nil {
		fmt.Println("error efd write: ", err)
		return
	}

	fmt.Println("write:  ", val)
*/
	val, err = efd.ReadEvents()
	if err !=nil {
		fmt.Println("error efd read: ", err)
		return
	}
	fmt.Println("read: ", val)


	/* TODO: register fd at kernel interface (for example cgroups memory watcher) */
	/* listen for new events */


	for {
		val, err := efd.ReadEvents()
		if err != nil {
			fmt.Printf("Error while reading from eventfd: %v", err)
			break
		}
		fmt.Printf("value: ", val)
	}


}

