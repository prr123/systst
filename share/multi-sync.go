package main

import (
	"fmt"
	"sync"
)

type pstatus struct {
	cs int
	iloop int
}

var stat [5]pstatus
var imax = 5

func worker(i int) {

// do process loop
	icount := 0
	for {
		icount++
// wait for action signal
		if stat[i].cs == 1 {
// execute the handler
			fmt.Println("proc ",i, " executing! ", icount,"loop: ", stat[i].iloop)
			stat[i].iloop++
			stat[i].cs = 0
		}
	} // end for loop
}


func main () {

var iproc int
var wg sync.WaitGroup

// initialise status
	for i:=0; i<imax; i++ {
		stat[i].cs = 0
		stat[i].iloop =0
	}


//create processes
	for i:=0; i<imax; i++ {
		wg.Add(1)
		go worker(i)
	}

// need an effective poll to see which process is free

	istart := 0

	for j:= 0; j< 10; j++ {
		fmt.Println("stat: ",stat)
		for it := 0; it < imax; it++ {
			iproc = (it + istart)%imax
// find the first free process
			if stat[iproc].cs == 0 {
				stat[iproc].cs = 1
			}
		}

	} // end for loop
}
