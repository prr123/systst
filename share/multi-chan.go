package main

import (
	"fmt"
	"time"
)

type pstatus struct {
	cs bool
	ch chan bool
}

var stat [5]pstatus
var imax = 5

var ready chan bool
var rstat bool

func worker(i int) {

// do process loop
//	icount := 0
	for {
// wait for action signal
		fmt.Println("proc ",i, " waiting! ")
		<- stat[i].ch
//		stat[i].cs = false
// execute the handler
		fmt.Println("proc ",i, " executing! ")
		time.Sleep(250 * time.Millisecond)
		stat[i].cs = true
		stat[i].ch <- false
		ready <- true
	} // end for loop
}


func main () {

var iproc int
var found bool

// initialise status
	for i:=0; i<imax; i++ {
		stat[i].cs = true
		stat[i].ch = make(chan bool)
//		stat[i].cs <- stat[i].ch
	}

	for i:=0; i<imax; i++ {
		fmt.Printf("%v ",stat[i].cs)
	}
	fmt.Println()

	ready = make(chan bool)

//create processes
	for i:=0; i<imax; i++ {
		fmt.Println("starting worker: ", i)
		go worker(i)
	}

// need an effective poll to see which process is free

	istart := 0
	

	for j:= 0; j< 10; j++ {
// new job
		fmt.Println("job: ",j) 
		found = false

		for i:=0; i<imax; i++ {
			fmt.Printf("%v ",stat[i].cs)
		}
		fmt.Println()

// see whether there is a free worker
		for it := 0; it < imax; it++ {

			iproc = (it + istart)%imax
// find the first free process
			if stat[iproc].cs {
// no race condition if changed before sending a signal
				fmt.Println("signal worker: ", iproc)
				stat[iproc].cs = false
// start the worker
				stat[iproc].ch <- true
				found = true
				break
			}

		} // end worker loop
		fmt.Println("selected worker: ", found)
		if !found {
			<- ready
// need to block until we receive a signal from one of the workers that it is free
			
		}

	} // end for loop
}
