package main


import (
	"fmt"
	"log"
	"os"
	"os/exec"
//	"github.com/erikdubbelboer/gspt"
	"strconv"
)


func main() {

	if len(os.Args) < 2 {
		fmt.Println("Need to specify number of clients!")
		fmt.Println("Usage: shm_start_clients numclients")
		os.Exit(1)
	}

	num_clients, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("error converting argument! ", err)
		fmt.Println("Usage: shm_start_clients numclients")
		os.Exit(1)
	}

    file, err := os.OpenFile("shm_clients.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    log.SetOutput(file)
    log.Println("Starting shm_clients")
/*
	path, err := os.Executable()
	if err != nil {
    	log.Println(err)
	}
	cl_nam := path + "shm_client"
	fmt.Println(path, cl_nam)
*/
	path2, err := os.Getwd()
	if err != nil {
    	log.Println(err)
	}
	cl_nam2 := path2 + "/shm_client"
//	fmt.Println(path2, cl_nam2)


	ab := byte(65)

	for i:=0; i< num_clients; i++ {
		character := string(ab)
//		fmt.Println(i, " : ", character)
		cmd := exec.Command(cl_nam2, character)
    	err := cmd.Start()
    	if err != nil {
			fmt.Println("client ",i,": ",err)
			log.Fatal(err)
    	}
		ab += 1
	} // end for loop
	fmt.Println("Started ", num_clients, " clients!")
	log.Println("Started ", num_clients, " clients!")
}
