package main

import (
	"fmt"
	"time"
)

// version 2: make ready a buffered integer channel

type pstatus struct {
	cs bool
	ch chan bool
}

var stat [5]pstatus
var imax = 5

var ready chan int

var proc_free_sig chan int
var rstat bool
var iact int




func worker(i int) {

	var recsig bool
// do process loop
//	icount := 0
	for {
// wait for action signal
		fmt.Println("proc ",i, " waiting! ")
		recsig = <-stat[i].ch

// execute the handler
		fmt.Println("proc ",i, " executing! ",recsig)
		time.Sleep(250 * time.Millisecond)
		stat[i].cs = true
//		stat[i].ch <- false
		ready <- i
	} // end for loop
}

func free_proc() {
	fmt.Println("free_proc: waiting for a free worker", iact)
	for iproc := range ready {
// need to block until we receive a signal from one of the workers that it is free
//new approach: check ready
		stat[iproc].cs = true
		iact--
		if (imax-iact) < 2 {
			fmt.Println("free_proc: signaling free worker: ", iproc)
			proc_free_sig <- iproc
		}
		fmt.Println("ready: worker ", iproc, " is ready! active: ", iact)

	}

}

func main () {


var iproc int

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

	ready = make(chan int, 5)
	proc_free_sig = make(chan int, 1)

//create processes
	for i:=0; i<imax; i++ {
		fmt.Println("starting worker: ", i)
		go worker(i)
	}

	go free_proc()
// need an effective poll to see which process is free

	istart := 0
	iact = 0
// job loop
	for j:= 0; j< 10; j++ {
// new job
		fmt.Println("job: ",j) 

		for i:=0; i<imax; i++ {
			fmt.Printf("%v ",stat[i].cs)
		}
		fmt.Println("active workers: ", iact)

		
// see whether there is a free worker
		if iact < imax {
			fmt.Println("worker available: ", iact)

			for it := 0; it < imax; it++ {

				iproc = (it + istart)%imax
// find the first free process
				if stat[iproc].cs {
// no race condition if changed before sending a signal
					fmt.Println("signal worker: ", iproc)
					stat[iproc].cs = false
// start the worker
					stat[iproc].ch <- true
					iact ++
					break
				}

			} // end worker loop
		} else {
// start one worker
			iproc = <-proc_free_sig
			fmt.Println("signal worker1: ", iproc)

			stat[iproc].cs = false
			stat[iproc].ch <- true
			iact++

		}

	} // end of job loop

// wait for jobs to complete
	close(proc_free_sig)

} // end main
