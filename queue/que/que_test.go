// ring buffer that stores msgs of arbitrary length
// method 1 store byte queue and add one byte as header with a length
// method 2 store structures with length and buffer


package rqueue


import (
	"testing"
	"fmt"
)

func Test_Init(t *testing.T) {
	var r *Ring
	var ring_size int = 300
	var buf_size int = 128
	var rdbuf, tstbuf []byte

	r = new(Ring)
	r.Init(ring_size)

	tstbuf = make([]byte,buf_size)
	for j:= 0; j<buf_size; j++ {
		tstbuf[j] = byte((j+1)%255)
	}
	err := r.Write(tstbuf)
	if err != nil {
		t.Error("Init: Expected nil, but received err")
	}
	rdbuf, err = r.Get()
	if err != nil {
		t.Error("Read: Expected nil, but received err")
	}

	for j:= 0; j<buf_size; j++ {
		if rdbuf[j] !=  tstbuf[j] {
			s := fmt.Sprint("Write Read: expected equal, but read unequal at position: ",j)
			t.Error(s)
			break
		}
	}

}

// test to see buffer overflow
func Test_Write_TM(t *testing.T){
	var r *Ring
	var ring_size int = 300
	var buf_size int = 128
	var tstbuf []byte

	r = new(Ring)
	r.Init(ring_size)

	tstbuf = make([]byte,buf_size)
	for j:= 0; j<buf_size; j++ {
		tstbuf[j] = byte((j+1)%255)
	}

	for i:= 0; i<4; i++ {
		err := r.Write(tstbuf)
		if i == 2 {
			if err == nil {
				t.Error("Write_tm: Expected error, but received nil")
			}
		} else {
			if err != nil {
				t.Error("Write_tm: Expected nil, but received err")
			}
		}
	} // for i loop
}

func Test_Read(t *testing.T){

}

/*
	var r *rqueue.Ring
	var tstbuf, readbuf []byte
	var buf_size int = 300

	fmt.Println("start")
	r = new(rqueue.Ring)
	r.Init(buf_size)
	fmt.Println("r: ", r)
	r.Disp()
	tstbuf = make([]byte,128)
	for j:= 0; j<len(tstbuf); j++ {
		tstbuf[j] = byte((j+1)%255)
	}
	fmt.Println("*****************************************************")
	fmt.Println("test buffer: ")
	fmt.Println(tstbuf)
	fmt.Println("*****************************************************")
	err := r.Write(tstbuf)
	if err != nil {
		fmt.Println("write ring error: ", err)
		return
	}
	r.Disp()
	fmt.Println("*****************************************************")
	fmt.Println("ring after write:")
	fmt.Println("*****************************************************")
	err = r.Write(tstbuf)
	r.Disp()
	fmt.Println("*****************************************************")
	fmt.Println("ring after write:")
	fmt.Println("*****************************************************")
	err = r.Write(tstbuf)
	r.Disp()
	
	readbuf, err = r.Get()
	fmt.Println("*****************************************************")
	fmt.Println("read buffer: ")
	fmt.Println(readbuf)
	fmt.Println("*****************************************************")

	r.Disp()
	
}
*/
