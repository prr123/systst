// ring buffer that stores msgs of arbitrary length
// method 1 store byte queue and add one byte as header with a length
// method 2 store structures with length and buffer


package main


import (
	"fmt"
	"queue/que"
)

func main() {

	var r *rqueue.Ring
	var tstbuf, readbuf []byte
	var buf_size int = 300
	var success bool
	var err error

	fmt.Println("start")
	r = new(rqueue.Ring)
	r.Init(buf_size)
	fmt.Println("r: ", r)
	r.Disp(false)
	tstbuf = make([]byte,128)
	for j:= 0; j<len(tstbuf); j++ {
		tstbuf[j] = byte((j+1)%255)
	}
	fmt.Println("*****************************************************")
	fmt.Println("test buffer: ")
	fmt.Println(tstbuf)
	fmt.Println("*****************************************************")
	n := r.Getcap()
	fmt.Println("ring cap: ", n)
	success, err = r.Write(tstbuf)
	if !success {
		fmt.Println("write: writing buf exceeds queue capacity!")
		return
	}
	if err != nil {
		fmt.Println("write ring error: ", err)
		return
	}
	fmt.Println("second write!")
	r.Disp(false)
	fmt.Println("*****************************************************")
	fmt.Println("ring after write:")
	fmt.Println("*****************************************************")
	success, err = r.Write(tstbuf)
	if !success {
		fmt.Println("write: writing buf exceeds queue capacity!")
		return
	}

	r.Disp(false)
	fmt.Println("*****************************************************")
	fmt.Println("ring after write:")
	fmt.Println("*****************************************************")
	fmt.Println("third write!")
	success, err = r.Write(tstbuf)
	if !success {
		fmt.Println("write: writing buf exceeds queue capacity!")
	}

	r.Disp(false)
	readbuf, err = r.Get()
	fmt.Println("*****************************************************")
	fmt.Println("read buffer: ")
	fmt.Println(readbuf)
	fmt.Println("*****************************************************")
	r.Disp(true)


	fmt.Println("fourth write")
	success, err = r.Write(tstbuf)
	if !success {
		fmt.Println("write: writing buf exceeds queue capacity!")
		return
	}
	if err != nil {
		fmt.Println("Write: ", err)
		return
	}

	fmt.Println("second read!")
	readbuf, err = r.Get()
	fmt.Println("*****************************************************")
	fmt.Println("read buffer: ")
	fmt.Println(readbuf)
	fmt.Println("*****************************************************")
	r.Disp(true)

	fmt.Println("third read!")
	readbuf, err = r.Get()
	fmt.Println("*****************************************************")
	fmt.Println("read buffer: ")
	fmt.Println(readbuf)
	fmt.Println("*****************************************************")
	r.Disp(true)

	r.Disp(true)
	fmt.Println("fifth write: should be an error")
	success, err = r.Write(tstbuf)
	if !success {
		fmt.Println("write: writing buf exceeds queue capacity!")
	}

//	num := r.Getcap()
//	fmt.Println("ring cap: ", num)

	if err != nil {
		fmt.Println("Write: ", err)
	}
	r.Disp(true)	

	

}

