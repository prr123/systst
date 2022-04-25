package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {

    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
//	fmt.Println("SIGIO: ", syscall.SIGIO)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGIO)

    go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
        done <- true
    }()

    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}
