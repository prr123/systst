package posixipc

// prr
// fixed func Swap

import "os"

// Proc is a wrapper around a generic process
type Proc struct {
	pid  int
	proc *os.Process
}

// ProcSlice is a slice of Lwps that can be sorted by pid
type ProcSlice []Proc

// Len returns the length of the ProcSlice
func (l ProcSlice) Len() int { return len(l) }

// Less returns either True or False, in regards to the size of a
// Lwp instance in a ProcSlice
func (l ProcSlice) Less(i, j int) (res bool) {
// l[i] < l[j]
	res = false
	if l[i].pid < l[j].pid {
		res = true
	}
return res
}

// Swap swaps the order of two Lwp objects in a ProcSlice
func (l ProcSlice) Swap(i, j int) {
	var tproc Proc
	tproc.pid = l[i].pid
	tproc.proc = l[i].proc
	l[i].pid = l[j].pid
	l[i].proc = l[j].proc
	l[j].pid = tproc.pid
	l[j].proc = tproc.proc
	return
}
