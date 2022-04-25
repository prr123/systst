package main

import (
	"fmt"
	"time"
	"sync"
)

// version 2: make ready a buffered integer channel


var ready chan int
var queue_jobs chan int

var proc_done chan bool


func worker(i int, wg *sync.WaitGroup) {
	var job int
// do process loop
//	icount := 0
	for {
// needs some thoughts
// wait for action signal
		fmt.Println("proc ",i, " waiting! ")
		job = <-queue_jobs
// execute the handler
		fmt.Println("proc ",i, " executing job: ", job)
		time.Sleep(250 * time.Millisecond)
//		ready <- i
// too early
		wg.Done()
	} // end for loop
}

func proc_end(rjob int) {
	fmt.Println("proc_end: starting! jobs: ", rjob)
	for {
		for iproc := range ready {
// need to block until we receive a signal from one of the workers that it is free
//new approach: check ready
			rjob--
			fmt.Println("proc-end: worker ", iproc, " is free! jobs remaining: ",rjob)
			if rjob == 0 {
				proc_done <- true
				break
			}
		} // iproc loop
	break
	} // for loop
}

func main () {

	var imax =5
	var numjobs = 10
	var wg sync.WaitGroup
// the ready buffer is limited to the number of workers
	ready = make(chan int, imax)
	proc_done = make(chan bool)
// test make job queue bigger than workers and less than numjobs
	queue_jobs = make(chan int, 8)

//create worker processes
	for i:=0; i<imax; i++ {
		fmt.Println("starting worker: ", i)
		go worker(i, &wg)
	}

//	go proc_end(numjobs)
// need an effective poll to see which process is free

// job loop
// start the initial batch of jobs with wnum of available workers
	for j:= 0; j< numjobs; j++ {
// start the worker
		wg.Add(1)
		queue_jobs <- j
	} // end of job loop

// wait for jobs to complete
	fmt.Println("main waiting for completion!")
//	<-proc_done
	wg.Wait()
	fmt.Println("finished")
} // end main
