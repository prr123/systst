package main

import (
	"fmt"
	"time"
)

// version 2: make ready a buffered integer channel

type pstatus struct {
	cs bool
	ch chan int
}

var stat [5]pstatus
var imax = 5

var ready chan int

var proc_free_sig chan int
var proc_done chan bool
var rstat bool
var iact int
var numjob int



func worker(i int) {

	var job int
// do process loop
//	icount := 0
	for {
// wait for action signal
		fmt.Println("proc ",i, " waiting! ")
		job = <-stat[i].ch
// execute the handler
		fmt.Println("proc ",i, " executing job: ", job)
		time.Sleep(250 * time.Millisecond)
		stat[i].cs = true
		ready <- i
	} // end for loop
}

func free_proc() {
	fmt.Println("free_proc: waiting for a free worker", iact)
	for {
//		if iact > 0 {
			for iproc := range ready {
// need to block until we receive a signal from one of the workers that it is free
//new approach: check ready
				stat[iproc].cs = true
				numjob--
				iact--
				if (imax-iact) < 2 {
					fmt.Println("free_proc: signaling free worker: ", iproc)
					proc_free_sig <- iproc
				}
				fmt.Println("ready: worker ", iproc, " is ready! active: ", iact, " jobs remaining: ",numjob)
			} // iproc loop
//		}

		if numjob == 0 {
				proc_done <- true
		}
	} // for loop

}

func main () {


var iproc int

// initialise status
	for i:=0; i<imax; i++ {
		stat[i].cs = true
		stat[i].ch = make(chan int)
//		stat[i].cs <- stat[i].ch
	}

	for i:=0; i<imax; i++ {
		fmt.Printf("%v ",stat[i].cs)
	}
	fmt.Println()

	ready = make(chan int, 5)
	proc_free_sig = make(chan int, 1)
	proc_done = make(chan bool)
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
	numjob =10
	for j:= 0; j< numjob; j++ {
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
					fmt.Println("signal worker: ", iproc, "job ", j)
					stat[iproc].cs = false
// start the worker
					stat[iproc].ch <- j
					iact ++
					break
				}

			} // end worker loop
		} else {
// start one worker
			iproc = <-proc_free_sig
			fmt.Println("signal worker1: ", iproc, " job ", j)

			stat[iproc].cs = false
			stat[iproc].ch <- j
			iact++

		}

	} // end of job loop
	fmt.Println("waiting for completion!")
// wait for jobs to complete
	<-proc_done
	fmt.Println("finished")
} // end main
