package main

import (
	"fmt"
//	"log"
	"time"
	"bufio"
	"os"
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
// delay until client has started
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    _, _ = reader.ReadString('\n')
    fmt.Println("starting to send ", maxTickNum, " messages!")

	count := 0
	for {
		count++
		mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0)
		fmt.Println("Sent message ", count)

		if count >= maxTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
