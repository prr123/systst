package posixipc

// corrections prr
/*
	func Randomint
   1. randInt replaced by randIntn (see https://golang.org/pkg/math/rand/)
   2. replace crypto/rand with math/rand
   3. added seed

	func generateMID
	inlined Randomint

	moved const MidMin and MidMax

	func NewMemoryDefaultHeap()
	- replaced newMID with generateMid
	destroy is not implemented

*/

import (
	"math/rand"
	"fmt"
	"sync"
	"time"
)


// DefaultHeapSize is the set size of a Memory object's heap,
// referenced when using AllocDefaultHeap().
const DefaultHeapSize = 256

// Default min and max values for new Mids
const (
	// Smallest allowed integer representing a Mid
	MidMin = 1
	// Larget allowed integer representing a Mid
	MidMax = 999
)



// Memory provides a segment of memory that can be shared
// between goroutines, passed via a mutex locker.
type Memory struct {
	mid int        // id
	mu  sync.Mutex // memory lock for sync
	desc MemDesc
	heap []byte
}

// MemDesc describes a shared piece of memory and is used
// as an identifier in the SharedMemory interface.
type MemDesc struct {
	mid int // id
	free   int // unused memory
	frames int // pages
}

// NewMemoryDefaultHeap creates a new Memory object that has a heap
// size of 256Kb, which is set as the default size when creating a
// new Memory object
func (m Memory) NewMemoryDefaultHeap() *Memory {
	heap := alloc(DefaultHeapSize)
//	mid := m.NewMid()
	mid := m.generateMid()
	return &Memory{
		mid:  mid,
		heap: heap,
	}
}

// NewMemory creates a Memory object with a heap of size N
func (m Memory) NewMemory(size int) *Memory {
	heap := alloc(size)
	mid := m.generateMid()
	return &Memory{
		mid:  mid,
		heap: heap,
	}
}

func alloc(size int) []byte {
	heap := make([]byte, size)
	return heap
}

func fillSlice(mem []byte) {
	_, err := rand.Read(mem)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}

func (m Memory) Dealloc(mid int) error {
	m.destroy(mid)
	return nil
}

func (m Memory) destroy(mid int) {
	// Release bytes created via Alloc()

}

// Mid is a unique id for a Memory type
type Mid struct {
	mid int
}

// NewMid creates a random integer that represents a Mid type
/*
func (m Memory) NewMid() *Mid {
	mid := m.generateMid()
	return &Mid{mid: mid}
}
*/
func (m Memory) NewMid() int {
	mid := m.generateMid()
	return mid
}

func (m Memory) generateMid() int {
	rand.Seed(time.Now().UnixNano())
	midInt := MidMin + rand.Intn(MidMax-MidMin)
	return midInt
}


// Returns an int >= min, < max
func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	return MidMin + rand.Intn(MidMax-MidMin)
}
