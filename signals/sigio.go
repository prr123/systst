package main

import (
	"fmt"
	"os"
	"os/signal"

	"syscall"

	"log"
)

func main() {
	var (
		tty   int
		err   error
		sigio = make(chan os.Signal)
	)
	tty, err = syscall.Open("/dev/tty", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		log.Fatalf("Cannot open tty port: %v\n", err)
	}

	signal.Notify(sigio, syscall.SIGIO)

	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(tty),
		uintptr(syscall.F_SETFL),
		uintptr(syscall.O_ASYNC|syscall.O_NONBLOCK))
	if errno != 0 {
		log.Fatalf("SYS_FCNTL call failed: %v\n", errno)
	}

	defer func() {
		syscall.Close(tty)
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		var buf = make([]byte, 8192)
		for {
			select {
			case <-done:
				log.Println("Quit!")
				return
			case <-sigio:
				log.Println("Read..")
				nr, err := syscall.Read(tty, buf)
				log.Println("Received:", nr, err)
				if err == syscall.EAGAIN {
					break
				}
				if nr > 0 {
					log.Printf("Write:%s(%d)\n", string(buf[:nr]), nr)
				}
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Signal rec: ",sig)
		close(done)
	}()

	fmt.Println("Ctrl+C to quit")
	<-done
	fmt.Println("exiting")
}
