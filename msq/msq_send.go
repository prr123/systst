package main

import (
	"fmt"
//	"log"
	"time"
	posix_mq "msq/msqlib"
)

const maxTickNum = 10

func main() {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mq, err := posix_mq.NewMessageQueue("test_queue", oflag, 0666, 6, 256)
	if err != nil {
		fmt.Println("error creating queue!", err)
		return
	}
	defer mq.Close()

	count := 0
	for {
		count++
		mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0)
		fmt.Println("Sent a new message")

		if count >= maxTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
