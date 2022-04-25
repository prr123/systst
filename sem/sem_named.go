package main

import (
	"fmt"
	sem "sem/semlib"
	"os"
)

func main() {
	var val uint32 = 4

// create a semaphore opbject
	nsem := new(sem.Semaphore)


	err := nsem.Open("new_sem", 0666, val)
	if err != nil {
		fmt.Println("error creating a  semaphore: ", err)
		os.Exit(1)
	}

	fmt.Println("semaphore created!")

}
