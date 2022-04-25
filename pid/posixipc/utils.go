package posixipc

// prr correction
// MaxProcs
//

import (
	"runtime"
)

type memstat_t runtime.MemStats

//NumCPU returns the number of CPUs
func NumCPU() int {
	return runtime.NumCPU()
}

// MaxProcs sets GOMAXPROCS to the number
// of CPUs reported by NumCPU() and returns the number
func MaxProcs() int {
	num := runtime.NumCPU()
	return runtime.GOMAXPROCS(num)
}

func MemoryStats() *runtime.MemStats {
	memstat := new(runtime.MemStats)
	runtime.ReadMemStats(memstat)
	return memstat
}

// SlicePtrFromStrings converts a slice of strings to a slice of
// pointers to NUL-terminated byte arrays. If any string contains
// a NUL byte, it returns (nil, EINVAL).
/*
func SlicePtrFromStrings(ss []string) ([][]byte, error) {
	var err error
	var bptr []byte
	bb := make(*[]byte, len(ss)+1)
	for i := 0; i < len(ss); i++ {
		bptr = []byte(ss[i])
		bb[i] = bptr
	}
	bb[len(ss)] = nil
	return bb, nil
}
*/
