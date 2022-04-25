package main

import (
	"fmt"
//	"log"

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

	count := 0
	for {
		count++

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
