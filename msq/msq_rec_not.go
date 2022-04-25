package main

import (
	"fmt"
//	"log"
	"syscall"
	"os/signal"
	"os"
	posix_mq "msq/msqlib"
)

const maxTickNum = 10

func main() {
	oflag := posix_mq.O_RDONLY
	mq, err := posix_mq.GetMessageQueue("test_queue", oflag, 0666, 6, 256)
	if err != nil {
		fmt.Println("error opening message queue: ",err)
		return
	}
	defer mq.Close()

	fmt.Println("Start receiving messages")

    c := make(chan os.Signal, 1)

// can specify as second arg which signal is redirected to channel c
    signal.Notify(c)

	count := 0
	for {
		count++

		mq.Notify(syscall.SIGUSR1)

    	fmt.Println("waiting for signal!")
    // Block until any signal is received.
    	s := <-c
    	fmt.Println("Got signal:", s)

		msg, _, err := mq.Receive()
		if err != nil {
			fmt.Println("error receive: ", err)
			return
		}
		fmt.Printf(string(msg))

		if count >= maxTickNum {
			break
		}
	}
}
