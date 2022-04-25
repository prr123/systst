// ring buffer that stores msgs of arbitrary length
// method 1 store byte queue and add one byte as header with a length
// method 2 store structures with length and buffer
// 26/9/2020
// add getb(buf) get with provided buffer
// modify Get() method: add read error if r.nmsgs == 0
// modify Write() add error msg if prvided buffer is larger than ring capacity
// add GetNmsgs(r *Ring)
// add getcap

package rqueue


import (
	"fmt"
	"errors"
	"encoding/binary"
)

type Ring struct {
	que []byte
	rp,wp, nb, nmsg, bs int
}

func (r *Ring) Init (size int) {
//	var ring Ring
	r.que = make([]byte, size)
	r.rp = 0
	r.wp = 0
	r.nb = 0
	r.nmsg = 0
	r.bs = size
	return
}

func (r *Ring) Init_wp (size int, wp int) {
//	var ring Ring
	r.que = make([]byte, size)
	r.rp = wp
	r.wp = wp
	r.nb = 0
	r.nmsg = 0
	r.bs = size
	return
}



func (r *Ring) Disp (iring bool) {

	fmt.Println("********************************************************")
	fmt.Println("Ring Buffer:")
	fmt.Println("Ring Size:       ",len(r.que))
	fmt.Println("Capacity:        ",cap(r.que))
	fmt.Println("Messages stored: ",r.nmsg)
	fmt.Println("Bytes stored:    ",r.nb)
	fmt.Println("Read position:   ",r.rp)
	fmt.Println("Write position:  ",r.wp)
	if iring {
		fmt.Println("Queue:")
		for i:=0; i< len(r.que); i++ {
			fmt.Print("| ",r.que[i])
		}
		fmt.Println()
	}
	fmt.Println("********************************************************")
}

func (r *Ring) Get_msg_num() (n int){
	n = r.nmsg
	return n
}

func (r *Ring) Getcap() (n int){

//	n = (r.bs + r.rp - r.wp)%r.bs

	if r.wp >= r.rp {
		n = r.bs - r.wp + r.rp
	} else {
		n = r.rp - r.wp
	}

	return n
}

func (r *Ring) Write (b []byte) (suc bool, err error) {
	var ilen16 int16
	var head [2]byte
	var ilen, icap, rcap, cp1 int

// first cond ring capacity
	ilen = len(b)
//	icap = (r.bs + r.rp - r.wp)%r.bs

	if r.wp >= r.rp {
		icap = r.bs - r.wp + r.rp
	} else {
		icap = r.rp - r.wp
	}

	if  (ilen +2) > icap {
		return false, errors.New("Ring Write: input buffer is larger than capacity!")
	}

	head_slice := head[0:2]
	ilen16 = int16(ilen)

	binary.LittleEndian.PutUint16(head_slice, uint16(ilen16))
//	fmt.Println("head sl: ", head_slice)
//	fmt.Println("write: len: ", ilen)

// second issue test whether writing needs to be wrapped around
	rcap = len(r.que) - r.wp
	switch true {
	case rcap > ilen+2:
// write header plus array on top of last msg
		copy(r.que[r.wp:],head_slice)
//	fmt.Println(" head: ", n)
		r.wp = r.wp + 2
		copy(r.que[r.wp:],b)
		r.wp = r.wp+ilen

	case rcap == ilen+2:
// write header plus array on top of last msg
		copy(r.que[r.wp:],head_slice)
//	fmt.Println(" head: ", n)
		r.wp = r.wp + 2
		copy(r.que[r.wp:],b)
		r.wp = 0

	case rcap == 2:
// write header on top of que and msg in bottom
		copy(r.que[r.wp:],head_slice)
		//r.wp = 0
		copy(r.que[r.wp:],b)
		r.wp = ilen

	case rcap == 1:
		copy(r.que[r.wp:],head_slice[0:1])
		copy(r.que[0:],head_slice[1:2])
		//r.wp = 1
		copy(r.que[r.wp:],b)
		r.wp = ilen + 1


	case rcap < ilen+2:
//
		cp1 = rcap-2
//		cp2 = ilen+2 - rcap
	 	copy(r.que[r.wp:],head_slice)
		r.wp = r.wp + 2
		copy(r.que[r.wp:],b[0:cp1])
		r.wp = 0
		copy(r.que[r.wp:],b[cp1:])
		r.wp = ilen - rcap  + 2

	}
	r.nmsg++
	r.nb = r.nb + ilen + 2
	return true, nil
}

func (r *Ring) Get() (buf []byte, err error) {
//	var ilen16 int16
	var head [2]byte
	var ilen, rcap, cp1 int
//first get the header

	if r.nmsg == 0 {
		err = errors.New("Ring: no messages stored!")
		return nil, err
	}

	rcap = len(r.que) - r.rp

	switch true {
	case rcap > 2:
 		head[0] = r.que[r.rp]
		r.rp++
		head[1] = r.que[r.rp]
		r.rp++

	case rcap == 2:
 		head[0] = r.que[r.rp]
		r.rp++
		head[1] = r.que[r.rp]
		r.rp = 0

	case rcap == 1:
 		head[0] = r.que[r.rp]
		r.rp = 0
		head[1] = r.que[r.rp]
		r.rp++

	default:
		s := fmt.Sprint("get error: rcap ", rcap)
		err= errors.New(s)
		return nil,err
	}

	head_slice := head[0:2]

	ilen = int(binary.LittleEndian.Uint16(head_slice))

//	fmt.Println("read head: ", r.rp, ilen)
	buf = make([]byte,ilen)
	rcap = len(r.que) - r.rp

	switch true {
// easy case no wrap around
	case rcap > ilen:
		cp1 = r.rp + ilen
		copy(buf, r.que[r.rp:cp1])
		r.rp = cp1

	case rcap == ilen:
		copy(buf, r.que[r.rp:])
		r.rp = 0

	case rcap < ilen:
// end of que
		cp1 = copy(buf, r.que[r.rp:])
// remainder is in beginning of que
		r.rp = ilen - cp1
		copy(buf[cp1:], r.que[:r.rp])
	}
	r.nmsg--
	r.nb = r.nb - ilen - 2
//	fmt.Println("get-copied: ", n, ",",cp1)
	return buf, nil
}

func (r *Ring) Getb(inbuf []byte) (outbuf []byte, err error) {
//	var ilen16 int16
	var head [2]byte
	var ilen, rcap, cp1 int
//first get the header

	if r.nmsg == 0 {
		err = errors.New("Ring: no messages stored!")
		return inbuf, err
	}

	rcap = len(r.que) - r.rp

	switch true {
	case rcap > 2:
 		head[0] = r.que[r.rp]
		r.rp++
		head[1] = r.que[r.rp]
		r.rp++

	case rcap == 2:
 		head[0] = r.que[r.rp]
		r.rp++
		head[1] = r.que[r.rp]
		r.rp = 0

	case rcap == 1:
 		head[0] = r.que[r.rp]
		r.rp = 0
		head[1] = r.que[r.rp]
		r.rp++

	default:
		s := fmt.Sprint("get error: rcap ", rcap)
		err= errors.New(s)
		return nil,err
	}

	head_slice := head[0:2]

	ilen = int(binary.LittleEndian.Uint16(head_slice))

//	fmt.Println("read head: ", r.rp, ilen)

	if len(inbuf) < ilen {
		s := fmt.Sprint("Ring-Getb: supplied buffer is too small! in: ", ilen)
		return inbuf, errors.New(s)
	}


	rcap = len(r.que) - r.rp

	switch true {
// easy case no wrap around
	case rcap > ilen:
		cp1 = r.rp + ilen
		copy(inbuf, r.que[r.rp:cp1])
		r.rp = cp1

	case rcap == ilen:
		copy(inbuf, r.que[r.rp:])
		r.rp = 0

	case rcap < ilen:
// end of que
		cp1 = copy(inbuf, r.que[r.rp:])
// remainder is in beginning of que
		r.rp = ilen - cp1
		copy(inbuf[cp1:], r.que[:r.rp])
	}
	r.nmsg--
	r.nb = r.nb - ilen - 2
//	fmt.Println("get-copied: ", n, ",",cp1)
	return inbuf, nil
}

